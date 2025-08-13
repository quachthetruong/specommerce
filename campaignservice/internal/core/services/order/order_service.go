// Package order implements iPhone campaign business logic.
//
// Campaign Rules:
// Gives away free M iPhones (configurable) for
// * Belong to first N customer (configurable) has successful transaction
// * Has a transaction greater than minimum (configurable)
// * One gift per customer (no duplicate)
// Orders processed in chronological order for fairness
//
// Two-Phase Processing Flow:
// Phase 1: ProcessPendingOrder - Validates and adds new orders to sorted set in chronological order
// Phase 2: ProcessOrderResult - Processes order completion and selects winners atomically
package order

import (
	"context"
	"fmt"
	"log"
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
	campaignRepo   secondary.CampaignRepository
	atomicExecutor atomicity.AtomicExecutor
	cacheClient    cache.Cache
	config         config.AppConfig
}

func NewOrderService(orderRepo secondary.OrderRepository, campaignRepo secondary.CampaignRepository, atomicExecutor atomicity.AtomicExecutor, cacheClient cache.Cache, config config.AppConfig) primary.OrderService {
	return &service{
		orderRepo:      orderRepo,
		campaignRepo:   campaignRepo,
		atomicExecutor: atomicExecutor,
		cacheClient:    cacheClient,
		config:         config,
	}
}

// ProcessPendingOrder adds a new order to the iPhone campaign pending orders sorted set.
// This is Phase 1 of the campaign flow. See ProcessOrderResult for Phase 2.
//
// The function executes an atomic Lua script that:
// 1. Validates order creation time is within campaign time window (start_time_micro to end_time_micro)
// 2. Early exits if campaign has reached maximum winners (policy_total_reward from campaign config)
// 3. Stores transaction data with PENDING status in Redis (transactions:{order_id})
// 4. Adds order to pending_orders sorted set with score = (created_at - start_time_micro)
//   - Score ensures chronological processing (earliest orders processed first)
//   - Relative scoring from campaign start time for consistent ordering
//
// Orders added here will be processed when their status changes from PENDING to
// SUCCESS/FAILED, maintaining order creation sequence for fair winner selection.
// All time values use microsecond precision for accurate chronological ordering.
func (s *service) ProcessPendingOrder(ctx context.Context, input order.Order) error {
	errTemplate := "orderService ProcessPendingOrder %w"
	luaScript := `
		local customer_id = KEYS[1]
		local created_at = tonumber(ARGV[1])
		local order_id = ARGV[2]
		local campaign_key = ARGV[3]
		local start_time_micro = tonumber(redis.call('HGET', campaign_key, 'start_time_micro'))
		local end_time_micro = tonumber(redis.call('HGET', campaign_key, 'end_time_micro'))
		local policy_total_reward = tonumber(redis.call('HGET',campaign_key, 'policy_total_reward')) or 0

		if created_at < start_time_micro or created_at > end_time_micro then
			return 
		end
		
		local winners_key = 'campaign_winners'
		local winner_count = redis.call('SCARD', winners_key)

		if winner_count == policy_total_reward then
			return 
		end

		local transaction_key = 'transactions:' .. order_id
		redis.call('HMSET', transaction_key, 'customer_id', customer_id, 'status', 'PENDING')
		
		local pending_orders_key = 'pending_orders'
		local score = created_at - start_time_micro
		redis.call('ZADD', pending_orders_key, score, order_id)
		return 
	`

	createdAt := input.CreatedAt.UnixMicro()
	campaignKey := fmt.Sprintf("campaign:%s", s.config.IphoneCampaign)
	_, err := s.cacheClient.Eval(ctx, luaScript, []string{input.CustomerId}, createdAt, input.Id.String(), campaignKey)
	if err != nil {
		return fmt.Errorf(errTemplate, err)
	}

	return nil
}

// ProcessOrderResult processes order completion for iPhone campaign winner selection.
// This is Phase 2 of the campaign flow. Processes orders added by ProcessPendingOrder.
//
// The function executes an atomic Lua script that:
// 1. Early exits if campaign has reached maximum winners (policy_total_reward)
// 2. Stores current transaction data (customer_id, status) in Redis
// 3. Updates customer's maximum transaction amount for successful orders
// 4. Checks if current order qualifies customer as immediate winner:
//   - Order status is SUCCESS
//   - Customer not already a winner
//   - Customer is in eligible set
//   - Customer's max amount >= minimum order amount policy
//
// 5. Recursively processes pending orders sorted set (oldest first):
//   - Stops immediately if encountering order still in PENDING status
//   - Skips failed orders and orders from existing winners (removes from sorted set and continues)
//   - Adds qualifying customers to eligible set (respects max tracked limit)
//   - Promotes eligible customers to winners if they meet amount threshold
//   - Continues until sorted set empty or winner quota reached
//
// Returns {has_new_winner, winners_count} indicating if any new winners were
// added and current total winner count. All operations are atomic to prevent
// race conditions in concurrent order processing.
func (s *service) ProcessOrderResult(ctx context.Context, order order.Order) error {
	errTemplate := "orderService ProcessOrderResult %w"
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
		local has_new_winner = false
		local is_campaign_finished = false

		local winners_count = redis.call('SCARD', winners_key)
		if winners_count == policy_total_reward then
			is_campaign_finished = true
			return {has_new_winner, is_campaign_finished}
		end

		local current_order_id_key = transaction_key .. ':' .. order_id
        local current_customer_id_key = customer_key .. ':' .. customer_id
        redis.call('HMSET', current_order_id_key, 'customer_id', customer_id, 'status', order_status)
        local current_max_total_amount = tonumber(redis.call('HGET', current_customer_id_key, 'max_total_amount')) or 0
        if order_status == 'SUCCESS' then
              current_max_total_amount = math.max(current_max_total_amount, order_total_amount)
			  redis.call('HSET', current_customer_id_key, 'max_total_amount', current_max_total_amount)
  		end
        
		if order_status == 'SUCCESS' and redis.call('SISMEMBER', winners_key, customer_id) == 0 and redis.call('SISMEMBER', eligible_key, customer_id) == 1 and current_max_total_amount >= policy_min_order_amount then
            redis.call('SADD', winners_key, customer_id)
			has_new_winner = true
		end

		local function recursive_pop() 
			local winners_count = redis.call('SCARD', winners_key)
            if winners_count == policy_total_reward then
				is_campaign_finished = true
                return {has_new_winner, is_campaign_finished}
			end

			local elements = redis.call('ZRANGE', pending_orders_key, 0, 0, 'WITHSCORES')
			
			if #elements == 0 then
				return {has_new_winner, is_campaign_finished}
			end
			
			local current_order_id = elements[1]
			local current_order_id_key = transaction_key .. ':' .. current_order_id
			local current_customer_id = redis.call('HGET', current_order_id_key, 'customer_id')
			local current_customer_id_key = customer_key .. ':' .. current_customer_id
			local current_status = redis.call('HGET', current_order_id_key, 'status')
            if current_status == 'PENDING' then
				return {has_new_winner, is_campaign_finished}
			end

			redis.call('ZREM', pending_orders_key, current_order_id)

			if current_status == 'FAILED' then
				return recursive_pop()
			end

            if redis.call('SISMEMBER', winners_key, current_customer_id) == 1 then
				return recursive_pop()  
			end

			if redis.call('SISMEMBER', eligible_key, current_customer_id) == 0 and redis.call('SCARD', eligible_key) == policy_max_tracked_orders then
				return recursive_pop()
			end

			redis.call('SADD', eligible_key, current_customer_id)

			local current_max_total_amount = tonumber(redis.call('HGET', current_customer_id_key, 'max_total_amount')) or 0


			if current_max_total_amount >= policy_min_order_amount then
				redis.call('SADD', winners_key, current_customer_id)
				has_new_winner = true
			end

			return recursive_pop()
		end
		return recursive_pop()  -- Start the recursion
	`
	campaignKey := fmt.Sprintf("campaign:%s", s.config.IphoneCampaign)
	result, err := s.cacheClient.Eval(ctx, luaScript, []string{order.CustomerId}, order.Id.String(), order.Status.String(), order.TotalAmount, campaignKey)
	if err != nil {
		return fmt.Errorf(errTemplate, err)
	}

	// Parse Lua script result: [has_new_winner, is_campaign_finished]
	if resultArray, ok := result.([]interface{}); ok && len(resultArray) == 2 {
		// Redis Lua returns booleans as integers: 1 = true, 0 = false
		hasNewWinner := false
		isCampaignFinished := false

		if val, ok := resultArray[0].(int64); ok && val == 1 {
			hasNewWinner = true
		}
		if val, ok := resultArray[1].(int64); ok && val == 1 {
			isCampaignFinished = true
		}

		log.Printf("Lua result: has_new_winner=%v, is_campaign_finished=%v", hasNewWinner, isCampaignFinished)

		if hasNewWinner && isCampaignFinished {
			campaign, err := s.campaignRepo.GetCampaignByType(ctx, "iphone")
			if err != nil {
				return fmt.Errorf(errTemplate, err)
			}

			// Get all winners from Redis campaign_winners set
			winners, err := s.cacheClient.SMembers(ctx, "campaign_winners")
			if err != nil {
				return fmt.Errorf(errTemplate, err)
			}

			log.Printf("Campaign finished! All winners: %v", winners)

			// Save all winners to database
			for _, customerID := range winners {
				err := s.campaignRepo.SaveWinner(ctx, campaign.Id, customerID)
				if err != nil {
					log.Printf("Failed to save winner %s: %v", customerID, err)
					// Continue with other winners even if one fails
				} else {
					log.Printf("Successfully saved winner: %s", customerID)
				}
			}
		}
	}

	return nil
}

func (s *service) SaveSuccessOrder(ctx context.Context, input order.Order) error {
	errTemplate := "orderService SaveSuccessOrder %w"
	_, err := s.orderRepo.Create(ctx, input)
	if err != nil {
		return fmt.Errorf(errTemplate, err)
	}
	return nil
}
