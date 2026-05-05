// app/dashboard/page.tsx
"use client";
import StatsCard from "@/components/Card";
import Navbar from "@/components/navbar";
import TransactionsList from "@/components/TransactionsList";
import Button from "@/components/Button";
import { useGetTransactions } from "@/hooks/useFetchTransactions";

export default function Dashboard() {
  const { transactions: recentTransactions } = useGetTransactions({ size: 5 });

  const { transactions: allTransactions, loading } = useGetTransactions({});

   const totalSavings = allTransactions
    .filter((t) => t.type === "saving")
    .reduce((sum, t) => sum + Number(t.amount), 0);

  const totalWithdrawals = allTransactions
    .filter((t) => t.type === "withdrawal")
    .reduce((sum, t) => sum + Number(t.amount), 0);

    const netBalance = totalSavings - totalWithdrawals;

  const fmt = (n: number) =>
    new Intl.NumberFormat("fr-CM", { maximumFractionDigits: 0 }).format(n) + " CFA";



  return (
    <div className="flex flex-col flex-1 min-h-screen bg-slate-50 dark:bg-slate-900">
      <Navbar />

      <main className="flex-1 w-full max-w-3xl mx-auto px-6 py-10 flex flex-col gap-8">
        <section className="flex gap-4">
         
          <StatsCard title="Total Savings" text={loading ? "..." : fmt(totalSavings)} />
          <StatsCard
            title="Total Withdrawals"
            text={loading ? "..." : fmt(totalWithdrawals)}
          />
        </section>


        {/* Net balance */}
        <section>
          {loading ? (
            <div className="p-6 bg-white dark:bg-slate-800 rounded-xl border border-slate-200 dark:border-slate-700 shadow-sm">
              <p className="text-sm font-medium text-slate-500 dark:text-slate-400 mb-1">Net Balance</p>
              <p className="text-3xl font-bold text-slate-800 dark:text-slate-100">...</p>
            </div>
          ) : (
            <div className={`p-6 rounded-xl border shadow-sm ${
              netBalance < 25
                ? "bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-800"
                : "bg-green-50 dark:bg-green-900/20 border-green-200 dark:border-green-800"
            }`}>
              <p className="text-sm font-medium text-slate-500 dark:text-slate-400 mb-1">Net Balance</p>
              <p className={`text-3xl font-bold ${
                netBalance < 25
                  ? "text-red-600 dark:text-red-400"
                  : "text-green-600 dark:text-green-400"
              }`}>
                {fmt(netBalance)}
              </p>
              <p className="text-xs mt-1 text-slate-400">
                {netBalance < 25
                  ? "⚠️ You can't buy anything under 25 CFA in Cameroon — top up your savings!"
                  : "You have enough to spend 🎉"}
              </p>
            </div>
          )}
        </section>
        <section className="flex gap-3">
          <Button text="Add Savings" onClick={() => { window.location.href = "/save"; }} />
          <Button text="Make Withdrawal" variant="secondary" onClick={() => { window.location.href = "/withdraw"; }} />
        </section>
        <section>
          <h2 className="text-lg font-semibold text-slate-700 dark:text-slate-300 mb-3">
            Recent Transactions
          </h2>
          <TransactionsList transactions={recentTransactions} />
        </section>
      </main>
    </div>
  );
}