import { gql } from "@apollo/client";

const ORDER_FIELDS = gql`
  fragment OrderFields on Order {
    id
    customerId
    restaurantId
    items {
      menuItemId
      name
      quantity
      unitPrice {
        amount
        currency
        display
      }
      specialInstructions
      subtotal {
        amount
        currency
        display
      }
    }
    status
    total {
      amount
      currency
      display
    }
    placedAt
    confirmedAt
    deliveredAt
    cancelledAt
    cancelReason
  }
`;

export const GET_ORDERS = gql`
  ${ORDER_FIELDS}
  query GetOrders($limit: Int, $offset: Int) {
    orders(limit: $limit, offset: $offset) {
      ...OrderFields
    }
  }
`;

export const GET_ORDER = gql`
  ${ORDER_FIELDS}
  query GetOrder($id: ID!) {
    order(id: $id) {
      ...OrderFields
    }
  }
`;

export const GET_ORDERS_BY_STATUS = gql`
  ${ORDER_FIELDS}
  query GetOrdersByStatus($status: OrderStatus!) {
    ordersByStatus(status: $status) {
      ...OrderFields
    }
  }
`;

export const PLACE_ORDER = gql`
  ${ORDER_FIELDS}
  mutation PlaceOrder($input: PlaceOrderInput!) {
    placeOrder(input: $input) {
      ...OrderFields
    }
  }
`;

export const ORDER_STATUS_CHANGED = gql`
  ${ORDER_FIELDS}
  subscription OrderStatusChanged($orderId: ID!) {
    orderStatusChanged(orderId: $orderId) {
      ...OrderFields
    }
  }
`;
