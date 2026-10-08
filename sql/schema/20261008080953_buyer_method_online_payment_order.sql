-- +goose Up
-- Drop the old, incorrect constraint
ALTER TABLE orders DROP CONSTRAINT IF EXISTS orders_method_check;

-- Add the corrected constraint
ALTER TABLE orders 
ADD CONSTRAINT orders_method_check 
CHECK (method IN ('cod', 'online'));

-- +goose Down
ALTER TABLE orders DROP CONSTRAINT IF EXISTS orders_method_check;
ALTER TABLE orders 
ADD CONSTRAINT orders_method_check 
CHECK (method IN ('cod', 'chapa'));