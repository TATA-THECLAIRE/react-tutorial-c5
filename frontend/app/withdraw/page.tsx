"use client";
import { saveTransaction } from "@/api/create-transaction";
import Button from "@/components/Button";
import Navbar from "@/components/navbar";
import { TransactionType } from "@/types/interfaces";
import { useRouter } from "next/navigation";
import React, { useState } from "react";
import { toast, ToastContainer } from "react-toastify";
import { useGetTransactions } from "@/hooks/useFetchTransactions";

function WithdrawPage() {
	const router = useRouter();
	const [amount, setAmount] = useState("");
	const [reason, setReason] = useState("");
	const [loading, setLoading] = useState(false);

	  // Fetch balance
	const { transactions } = useGetTransactions({});
	const totalSavings = transactions
		.filter((t) => t.type === "saving")
		.reduce((sum, t) => sum + Number(t.amount), 0);
	const totalWithdrawals = transactions
		.filter((t) => t.type === "withdrawal")
		.reduce((sum, t) => sum + Number(t.amount), 0);
	const netBalance = totalSavings - totalWithdrawals;

	const handleWithdraw = async() => {
		const withdrawAmount = Number(amount);
		if (!amount || isNaN(Number(amount))) {
			toast.error("Please enter a valid amount.");
			return;
		}
		if (withdrawAmount > netBalance) {
			toast.error(`Insufficient balance! You only have ${netBalance.toLocaleString("fr-CM")} CFA.`);
			return;
		}

		if (netBalance - withdrawAmount < 25) {
			toast.error("This withdrawal would leave you with less than 25 CFA — you can't buy anything under 25 CFA in Cameroon!");
			return;
		}
		setLoading(true);
		
		const payload: TransactionType = {
			amount: Number(amount),
			reason,
			type: "withdrawal",
			created_at: Date(),
		};
		const res = await saveTransaction(payload);
		setLoading(false);


		if (res.success) {
			toast("Withdrawal successful! Redirecting...");
			setTimeout(() => router.push("/dashboard"), 1500);
		} else {
			toast.error("Failed to withdraw money! Try again.");
		}
	};

	return (
		<div className="flex flex-col flex-1 min-h-screen bg-slate-50 dark:bg-slate-900">
			<Navbar />
			<main className="flex flex-1 items-center justify-center px-6 py-10">
				<div className="w-full max-w-md bg-white dark:bg-slate-800 rounded-xl border border-slate-200 dark:border-slate-700 shadow-sm p-8">
					<h1 className="text-xl font-semibold text-slate-800 dark:text-slate-100 mb-6">
						Make a Withdrawal
					</h1>
					          {/* Live balance display */}
					<p className={`text-sm mb-6 font-medium ${
						netBalance < 25 ? "text-red-500" : "text-slate-400"
					    }`}>
						Available balance:{" "}
						<span className="font-bold">
						{netBalance.toLocaleString("fr-CM")} CFA
						</span>
						{netBalance < 25 && " ⚠️ Too low to spend"}
					</p>
					<form className="flex flex-col gap-4">
						<div className="flex flex-col gap-1">
							<label className="text-sm font-medium text-slate-600 dark:text-slate-400">
								Amount
							</label>
							<input
							    type="number"
                                min="1"
								className="border border-slate-200 dark:border-slate-600 bg-white dark:bg-slate-900 rounded-lg px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:text-slate-100"
								placeholder="e.g. 1000"
								value={amount}
								onChange={(v) => setAmount(v.target.value)}
							/>
						</div>
						<div className="flex flex-col gap-1">
							<label className="text-sm font-medium text-slate-600 dark:text-slate-400">
								Reason
							</label>
							<input
								className="border border-slate-200 dark:border-slate-600 bg-white dark:bg-slate-900 rounded-lg px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:text-slate-100"
								placeholder="e.g. Groceries"
								value={reason}
								onChange={(v) => setReason(v.target.value)}
							/>
						</div>
						<div className="pt-2">
							<Button text="Withdraw" onClick={handleWithdraw} />
						</div>
					</form>
				</div>
			</main>
			<ToastContainer />
		</div>
	);
}

export default WithdrawPage;
