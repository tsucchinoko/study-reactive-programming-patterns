import React from "react";
import { ApolloProvider } from "@apollo/client/react";
import { NavigationContainer } from "@react-navigation/native";
import { StatusBar } from "expo-status-bar";
import { apolloClient } from "./src/shared/graphql/client";
import { AppNavigator } from "./src/navigation/AppNavigator";

export default function App() {
  return (
    <ApolloProvider client={apolloClient}>
      <NavigationContainer>
        <AppNavigator />
        <StatusBar style="auto" />
      </NavigationContainer>
    </ApolloProvider>
  );
}
