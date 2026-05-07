import { authHeaders } from "@/app/lib/auth";

const BASEURL = "http://localhost:8080";

export interface DashboardStats {
  balance: number;
  total_savings: number;
  total_withdrawals: number;
}

export const getDashboardStats = async (): Promise<DashboardStats | null> => {
  try {
    const res = await fetch(`${BASEURL}/api/v1/dashboard/stats`, {
      headers: authHeaders(),
    });
    if (!res.ok) return null;
    return await res.json();
  } catch {
    return null;
  }
};