import {
  ApolloClient,
  InMemoryCache,
  HttpLink,
  split,
} from "@apollo/client";
import { GraphQLWsLink } from "@apollo/client/link/subscriptions";
import { getMainDefinition } from "@apollo/client/utilities";
import { createClient } from "graphql-ws";
import { Platform } from "react-native";

// Android emulator uses 10.0.2.2 to reach host machine's localhost.
const HOST = Platform.OS === "android" ? "10.0.2.2" : "localhost";
const API_URL = `http://${HOST}:8080/graphql`;
const WS_URL = `ws://${HOST}:8080/graphql`;

const httpLink = new HttpLink({ uri: API_URL });

const wsLink = new GraphQLWsLink(
  createClient({
    url: WS_URL,
    retryAttempts: 5,
    shouldRetry: () => true,
  })
);

// Route subscriptions over WebSocket, everything else over HTTP.
const splitLink = split(
  ({ query }) => {
    const definition = getMainDefinition(query);
    return (
      definition.kind === "OperationDefinition" &&
      definition.operation === "subscription"
    );
  },
  wsLink,
  httpLink
);

export const apolloClient = new ApolloClient({
  link: splitLink,
  cache: new InMemoryCache(),
});
