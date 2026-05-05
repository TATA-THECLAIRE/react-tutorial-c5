import { GetTransactionsParamsType, TransactionType } from "@/types/interfaces";
import { authHeaders } from "@/app/lib/auth";

export interface GetTransactionsRes {
	transactions: TransactionType[];
	error: any;
}

const BASEURL = "http://localhost:8080";

export const getAllTransactions = async (
	query: GetTransactionsParamsType,
): Promise<GetTransactionsRes> => {
	const queries = new URLSearchParams ();
	if (query.size) {
		queries.set("size", String(query.size));
	}

	if (query.type) {
		queries.set("type", query.type);
	}

	const url =
		`${BASEURL}/api/v1/transactions${queries.toString() ? "?" + queries : ""}`;

	try {
		const res = await fetch(url, { headers: authHeaders() });
		if (!res.ok) return { transactions: [], error: "Failed to fetch" };

		const data = await res.json();
		// backend returns array directly
		const transactions = Array.isArray(data) ? data : [];
		return {
			transactions, error: null,
		};
	} catch (error) {
		return {
			transactions: [],
			error,
		};
	}
};



