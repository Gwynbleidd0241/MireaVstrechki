import { apiRequest } from "./client";

export type GeocodeResult = {
  name: string;
  description: string;
  coordinates: {
    latitude: number;
    longitude: number;
  };
  map_url: string;
};

export function geocode(address: string) {
  return apiRequest<GeocodeResult>(
    `/geo/geocode?address=${encodeURIComponent(address)}`,
  );
}
