// ---------------------------------------------------------------------------
// User & Auth
// ---------------------------------------------------------------------------

export type UserRole = 'customer' | 'owner';

export interface User {
  id: number;
  created_at: string;
  name: string;
  email: string;
  phone: string;
  role: UserRole;
  activated: boolean;
  version: number;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  authentication_token: string;
}

export interface RegisterRequest {
  name: string;
  email: string;
  password: string;
  phone: string;
  role: UserRole;
}

export interface RegisterResponse {
  user: User;
}

export interface ActivateRequest {
  token: string;
}

export interface ActivateResponse {
  user: User;
}

// ---------------------------------------------------------------------------
// Restaurant
// ---------------------------------------------------------------------------

export interface Restaurant {
  id: number;
  created_at: string;
  name: string;
  email: string;
  phone: string;
  street_address: string;
  opening_time: string;
  closing_time: string;
  delivery_fee: number;
  delivery_radius: number;
  is_open: boolean;
  version: number;
}

// ---------------------------------------------------------------------------
// Category
// ---------------------------------------------------------------------------

export interface Category {
  id: number;
  created_at: string;
  name: string;
  display_order: number;
  version: number;
}

// ---------------------------------------------------------------------------
// Dish
// ---------------------------------------------------------------------------

export interface Dish {
  id: number;
  created_at: string;
  name: string;
  description: string;
  price: number; // in smallest currency unit (paise/cents)
  category_id: number;
  image_url: string | null;
  is_available: boolean;
  is_vegetarian: boolean;
  is_featured: boolean;
  preparation_time: number;
  version: number;
}

// ---------------------------------------------------------------------------
// Address
// ---------------------------------------------------------------------------

export interface Address {
  id: number;
  created_at: string;
  customer_id: number;
  street_line_1: string;
  street_line_2?: string;
  city: string;
  state: string;
  postal_code: string;
  country: string;
  is_default: boolean;
  version: number;
}

export interface CreateAddressRequest {
  street_line1: string;
  street_line2?: string;
  city: string;
  state: string;
  postal_code: string;
  country: string;
  is_default: boolean;
}

export interface UpdateAddressRequest {
  street_line_1?: string;
  street_line_2?: string;
  city?: string;
  state?: string;
  postal_code?: string;
  country?: string;
  is_default?: boolean;
}

// ---------------------------------------------------------------------------
// Order
// ---------------------------------------------------------------------------

export type OrderStatus =
  | 'pending'
  | 'accepted'
  | 'preparing'
  | 'ready'
  | 'out_for_delivery'
  | 'delivered'
  | 'cancelled';

export type OrderType = 'delivery' | 'pickup';

export interface Order {
  id: number;
  created_at: string;
  customer_id: number;
  status: OrderStatus;
  order_type: OrderType;
  delivery_address_id?: number;
  pickup_time?: string;
  subtotal: number;
  delivery_fee: number;
  tax: number;
  total: number;
  version: number;
}

export interface OrderItem {
  id: number;
  created_at: string;
  order_id: number;
  dish_id: number;
  quantity: number;
  price_at_purchase: number;
  subtotal: number;
  version: number;
}

export interface CreateOrderItemRequest {
  dish_id: number;
  quantity: number;
}

export interface CreateOrderRequest {
  order_type: OrderType;
  delivery_address_id?: number;
  pickup_time?: string;
  items: CreateOrderItemRequest[];
}

export interface CreateOrderResponse {
  order: Order;
  items: OrderItem[];
}

export interface GetOrderResponse {
  order: Order;
  items: OrderItem[];
}

// ---------------------------------------------------------------------------
// Dashboard
// ---------------------------------------------------------------------------

export interface DashboardStats {
  restaurant_name: string;
  category_count: number;
  dish_count: number;
  pending_orders: number;
  accepted_orders: number;
  preparing_orders: number;
  ready_orders: number;
  out_for_delivery_orders: number;
  delivered_orders: number;
  cancelled_orders: number;
}

// ---------------------------------------------------------------------------
// API Envelope
// ---------------------------------------------------------------------------

export interface ApiError {
  error: string | Record<string, string>;
}