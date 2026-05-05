import { getAllTransactions } from "@/api/get-transactions";
import { GetTransactionsParamsType, TransactionType } from "@/types/interfaces";
import { useCallback, useEffect, useState } from "react";

export const useGetTransactions = (query: GetTransactionsParamsType) => {
	const [transactions, setTransactions] = useState<TransactionType[]>([]);
	  const [loading, setLoading] = useState(true);
	  const [error, setError] = useState<string | null>(null);

	const fetchData = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await getAllTransactions(query);
	  // backend returns array directly, not { transactions: [] }
	  const list = Array.isArray(res.transactions) ? res.transactions : [];
      setTransactions(list);
    } catch (e: any) {
      setError(e.message);
	  setTransactions([]);
    } finally {
      setLoading(false);
    }
  }, [query.size, query.type]);

  useEffect(() => { fetchData(); }, [fetchData]);

  return { transactions, loading, error, refetch: fetchData };
};