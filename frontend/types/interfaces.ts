export interface TransactionType {
	amount: number ;
	reason: string;
	created_at: string;
	type: "saving" | "withdrawal";
	id?: number;
}

export interface GetTransactionsParamsType {
	type?: "saving" | "withdrawal";
	size?: number;
}

export interface TransactionsResponse {
  transactions: TransactionType[];
}

export interface CreateTransactionResponse {
  success: boolean;
  error?: string;
}
