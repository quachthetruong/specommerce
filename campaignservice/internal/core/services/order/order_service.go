package order

import (
	"context"
	"fmt"
	"specommerce/campaignservice/config"
	"specommerce/campaignservice/internal/core/domain/order"
	"specommerce/campaignservice/internal/core/ports/primary"
	"specommerce/campaignservice/internal/core/ports/secondary"
	"specommerce/campaignservice/pkg/atomicity"
	"specommerce/campaignservice/pkg/cache"
)

// OrderService implements the order business logic
type service struct {
	orderRepo      secondary.OrderRepository
	atomicExecutor atomicity.AtomicExecutor
	cacheClient    cache.Cache
	config         config.AppConfig
}

func NewOrderService(orderRepo secondary.OrderRepository, atomicExecutor atomicity.AtomicExecutor, cacheClient cache.Cache, config config.AppConfig) primary.OrderService {
	return &service{
		orderRepo:      orderRepo,
		atomicExecutor: atomicExecutor,
		cacheClient:    cacheClient,
		config:         config,
	}
}

// ProcessPendingOrder processes a new order by saving it to the database and adding it to a Redis sorted set
// Step 1: Save the order to the database
// Step 2: Add the order to a Redis sorted set with a score based on the order creation time
func (s *service) ProcessPendingOrder(ctx context.Context, input order.Order) (int64, error) {
	errTemplate := "orderService CreateOrder %w"
	luaScript := `
		local customer_id = KEYS[1]
		local order_time = tonumber(ARGV[1])
		local order_id = ARGV[2]
		local campaign_key = ARGV[3]
		local start_time_micro = tonumber(redis.call('HGET', campaign_key, 'start_time_micro'))
		local end_time_micro = tonumber(redis.call('HGET', campaign_key, 'end_time_micro'))

		if order_time < start_time_micro or order_time > end_time_micro then
			return 0
		end
		
		local winners_key = 'campaign_winners'
		local winner_count = redis.call('SCARD', winners_key)

		if winner_count == 20 then
			return 
		end

		local transaction_key = 'transactions:' .. order_id
		redis.call('HMSET', transaction_key, 'customer_id', customer_id, 'status', 'PENDING')
		
		local pending_orders_key = 'pending_orders'
		local score = order_time - start_time_micro
		return redis.call('ZADD', pending_orders_key, score, order_id)
	`

	orderTime := input.CreatedAt.UnixMicro()
	campaignKey := fmt.Sprintf("campaign:%s", s.config.IphoneCampaign)
	pendingCount, err := s.cacheClient.Eval(ctx, luaScript, []string{input.CustomerId}, orderTime, input.Id.String(), campaignKey)
	if err != nil {
		return 0, fmt.Errorf(errTemplate, err)
	}
	if pendingCount == nil {
		return 0, fmt.Errorf("orderService ProcessPendingOrder: pending_count is nil")
	}
	return pendingCount.(int64), nil
}

// ProcessOrderResult processes the result of an order after it has been completed
// Goal: The system supports a hot-sale campaign that gives away 20 free iPhones for transactions matching requirements:
// - First 50 customers with successful transactions
// - Order amount greater than SGD 200 (configurable)
// - One gift per customer (no duplicate wins)

// Steps:
// 1. Check if winner count is less than 20
// 2. Update the transaction status in Redis
// 3. Update the customer's maximum transaction amount in Redis
func (s *service) ProcessOrderResult(ctx context.Context, order order.Order) (int64, error) {
	luaScript := `
		local customer_id = KEYS[1]
		local order_id = ARGV[1]
		local order_status = ARGV[2]
		local order_total_amount = tonumber(ARGV[3])
		local campaign_key = ARGV[4]
		local winners_key = 'campaign_winners'
		local eligible_key = 'campaign_eligible'
        local pending_orders_key = 'pending_orders'
		local transaction_key = 'transactions'
		local customer_key = 'customers'
		local policy_total_reward = tonumber(redis.call('HGET',campaign_key, 'policy_total_reward')) or 0
		local policy_min_order_amount = tonumber(redis.call('HGET', campaign_key, 'policy_min_order_amount')) or 0
		local policy_max_tracked_orders = tonumber(redis.call('HGET', campaign_key, 'policy_max_tracked_orders')) or 0

		local winners_count = redis.call('SCARD', winners_key)
		if winners_count == policy_total_reward then
			return winners_count
		end
        
 		local current_transaction_key = transaction_key .. ':' .. order_id
        redis.call('HMSET', current_transaction_key, 'customer_id', customer_id, 'status', order_status)
		local current_customer_key = customer_key .. ':' .. customer_id
        local current_max_total_amount = tonumber(redis.call('HGET', current_customer_key, 'max_total_amount')) or 0
        if order_status == 'SUCCESS' then
              current_max_total_amount = math.max(current_max_total_amount, order_total_amount)
			  redis.call('HSET', current_customer_key, 'max_total_amount', current_max_total_amount)
  		end
        
		if order_status == 'SUCCESS' and redis.call('SISMEMBER', winners_key, customer_id) == 0 and redis.call('SISMEMBER', eligible_key, customer_id) == 1 and current_max_total_amount > policy_min_order_amount then
            redis.call('SADD', winners_key, customer_id)
			winners_count = winners_count + 1
			return winners_count
		end

		local function recursive_pop()
			local elements = redis.call('ZRANGE', pending_orders_key, 0, 0, 'WITHSCORES')
			
			if #elements == 0 then
				return winners_count
			end
			
			local member = elements[1]
			local score = tonumber(elements[2])
 			local current_transaction_key = transaction_key .. ':' .. member
			local current_customer = redis.call('HGET', current_transaction_key, 'customer_id')
			local current_status = redis.call('HGET', current_transaction_key, 'status')
			redis.log(redis.LOG_NOTICE, "Processing order: " .. member .. ", status: " .. current_status)
            if current_status == 'PENDING' then
				return winners_count
			end

			redis.call('ZREM', pending_orders_key, member)

			if current_status == 'FAILED' then
				return recursive_pop()
			end

            if redis.call('SISMEMBER', winners_key, current_customer) == 1 then
				return recursive_pop()  
			end

			if redis.call('SISMEMBER', eligible_key, current_customer) == 0 and redis.call('SCARD', eligible_key) == policy_max_tracked_orders then
				return recursive_pop()
			end

			redis.call('SADD', eligible_key, current_customer)

			local current_customer = customer_key .. ':' .. current_customer
			local current_max_total_amount = tonumber(redis.call('HGET', current_customer, 'max_total_amount')) or 0

			if current_max_total_amount > policy_min_order_amount then
				redis.call('SADD', winners_key, current_customer)
 
            if redis.call('SCARD', winners_key) < policy_total_reward then
                return recursive_pop()
			end

			winners_count = redis.call('SCARD', winners_key)
			return winners_count
		end
		return recursive_pop()  -- Start the recursion
	`
	campaignKey := fmt.Sprintf("campaign:%s", s.config.IphoneCampaign)
	winnerCount, err := s.cacheClient.Eval(ctx, luaScript, []string{order.CustomerId}, order.Id.String(), order.Status.String(), order.TotalAmount, campaignKey)
	if err != nil {
		return 0, fmt.Errorf("orderService ProcessOrderResult %w", err)
	}
	if winnerCount == nil {
		return 0, fmt.Errorf("orderService ProcessOrderResult: winnerCount is nil")
	}
	return winnerCount.(int64), nil
}
