import React from "react";
import { StyleSheet, Text, View } from "react-native";

type Restaurant = {
  id: string;
  name: string;
  cuisine: string;
  location: { address: string };
  isOpen: boolean;
  menu: { items: Array<{ id: string }> };
};

type Props = {
  restaurant: Restaurant;
};

export function RestaurantCard({ restaurant }: Props) {
  return (
    <View style={styles.card}>
      <View style={styles.header}>
        <Text style={styles.name}>{restaurant.name}</Text>
        <View
          style={[
            styles.statusDot,
            { backgroundColor: restaurant.isOpen ? "#10B981" : "#EF4444" },
          ]}
        />
      </View>
      <Text style={styles.cuisine}>{restaurant.cuisine}</Text>
      <Text style={styles.address}>{restaurant.location.address}</Text>
      <Text style={styles.menuCount}>
        メニュー {restaurant.menu.items.length}品
      </Text>
    </View>
  );
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: "#fff",
    borderRadius: 12,
    padding: 16,
    marginHorizontal: 16,
    marginVertical: 6,
    shadowColor: "#000",
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.08,
    shadowRadius: 4,
    elevation: 2,
  },
  header: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    marginBottom: 4,
  },
  name: {
    fontSize: 16,
    fontWeight: "bold",
    color: "#1F2937",
  },
  statusDot: {
    width: 8,
    height: 8,
    borderRadius: 4,
  },
  cuisine: {
    fontSize: 13,
    color: "#3B82F6",
    marginBottom: 4,
  },
  address: {
    fontSize: 12,
    color: "#9CA3AF",
    marginBottom: 4,
  },
  menuCount: {
    fontSize: 12,
    color: "#6B7280",
  },
});
