import {
  ApolloClient,
  ApolloLink,
  InMemoryCache,
  HttpLink,
} from "@apollo/client";
import { GraphQLWsLink } from "@apollo/client/link/subscriptions";
import { getMainDefinition } from "@apollo/client/utilities";
import { createClient } from "graphql-ws";
import { Platform } from "react-native";

// Androidエミュレータはホストマシンのlocalhostに接続するために10.0.2.2を使用する。
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

// SubscriptionはWebSocket経由、それ以外はHTTP経由でルーティングする。
const splitLink = ApolloLink.split(
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
