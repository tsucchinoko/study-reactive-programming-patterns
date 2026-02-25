import { gql } from "@apollo/client";

export const GET_DELIVERY_BY_ORDER = gql`
  query GetDeliveryByOrder($orderId: ID!) {
    deliveryByOrder(orderId: $orderId) {
      id
      orderId
      driverId
      driver {
        id
        name
        phone
        status
        currentLocation {
          lat
          lng
        }
      }
      status
      assignedAt
      pickedUpAt
      deliveredAt
    }
  }
`;

export const DRIVER_LOCATION_UPDATED = gql`
  subscription DriverLocationUpdated($driverId: ID!) {
    driverLocationUpdated(driverId: $driverId) {
      driverId
      latitude
      longitude
      timestamp
    }
  }
`;
