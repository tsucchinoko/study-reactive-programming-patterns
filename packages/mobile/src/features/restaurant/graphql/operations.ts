import { gql } from "@apollo/client";

export const GET_RESTAURANTS = gql`
  query GetRestaurants {
    restaurants {
      id
      name
      cuisine
      location {
        address
      }
      menu {
        id
        items {
          id
          name
          category
          price {
            display
          }
          available
        }
      }
      isOpen
    }
  }
`;
