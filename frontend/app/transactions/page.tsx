import React, { Suspense } from "react";
import TransactionsContent from "./TransactionsContent";

export default function TransactionsPage() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <TransactionsContent />
    </Suspense>
  );
}