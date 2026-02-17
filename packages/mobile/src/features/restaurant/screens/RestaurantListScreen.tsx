import { useQuery } from "@apollo/client/react";
import React from "react";
import {
  ActivityIndicator,
  FlatList,
  StyleSheet,
  Text,
  View,
} from "react-native";
import { RestaurantCard } from "../components/RestaurantCard";
import { GET_RESTAURANTS } from "../graphql/operations";

export function RestaurantListScreen() {
  const { data, loading, error, refetch } = useQuery<{ restaurants: any[] }>(
    GET_RESTAURANTS,
  );

  if (loading && !data) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color="#3B82F6" />
      </View>
    );
  }

  if (error) {
    return (
      <View style={styles.center}>
        <Text style={styles.error}>エラー: {error.message}</Text>
      </View>
    );
  }

  const restaurants = data?.restaurants ?? [];

  return (
    <FlatList
      data={restaurants}
      keyExtractor={(item) => item.id}
      renderItem={({ item }) => <RestaurantCard restaurant={item} />}
      contentContainerStyle={styles.list}
      onRefresh={() => refetch()}
      refreshing={loading}
      ListEmptyComponent={
        <View style={styles.center}>
          <Text style={styles.empty}>レストランがありません</Text>
        </View>
      }
    />
  );
}

const styles = StyleSheet.create({
  list: {
    paddingVertical: 8,
  },
  center: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
    padding: 32,
  },
  error: {
    color: "#EF4444",
    fontSize: 14,
  },
  empty: {
    fontSize: 16,
    fontWeight: "600",
    color: "#6B7280",
  },
});
