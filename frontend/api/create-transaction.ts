import { CreateTransactionResponse, TransactionType } from "@/types/interfaces";
import { authHeaders } from "@/app/lib/auth"




const BASEURL = "http://localhost:8080";

export const saveTransaction = async (
  payload:{
	amount: number;
	reason: string;
	type: "saving" | "withdrawal";
  }
): Promise<CreateTransactionResponse> => {
  try {
    const res = await fetch(BASEURL + "/api/v1/transactions", {
      method: "POST",
      headers: authHeaders(),
      body: JSON.stringify(payload),
    });

    if (!res.ok) {
      const data = await res.json().catch(()=> ({}));
      return { success: false, error: data.error || "Request failed" };
    }

    return { success: true };
  } catch (err: any) {
    return { success: false, error: err.message };
  }
};