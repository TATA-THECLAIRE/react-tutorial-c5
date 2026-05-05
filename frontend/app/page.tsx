"use client";

import Link from "next/link";
import { useEffect, useState } from "react";

export default function WelcomePage() {
  const [visible, setVisible] = useState(false);
  useEffect(() => { setVisible(true); }, []);

  return (
    <div className="min-h-screen dark:bg-slate-900 flex flex-col">

      {/* Header */}
      <header className="flex items-center justify-between px-8 py-6">
        <span className="text-2xl font-black tracking-tighter text-white font-mono">
          piggy<span className="text-emerald-400">.</span>
        </span>
        <div className="flex items-center gap-3">
          <Link
            href="/login"
            className="text-sm text-zinc-400 hover:text-white transition font-medium px-4 py-2"
          >
            Sign in
          </Link>
          <Link
            href="/register"
            className="text-sm bg-emerald-500 hover:bg-emerald-400 text-zinc-950 font-bold px-4 py-2 rounded-lg transition"
          >
            Get started
          </Link>
        </div>
      </header>

      {/* Hero */}
      <main className="flex-1 flex flex-col items-center justify-center px-6 text-center">
        <div
          className={`transition-all duration-700 ${
            visible ? "opacity-100 translate-y-0" : "opacity-0 translate-y-6"
          }`}
        >
          {/* Piggy icon */}
          <div className="mb-8 flex justify-center">
            <div className="w-24 h-24 rounded-full bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-center text-5xl">
              🐷
            </div>
          </div>

          <h1 className="text-5xl sm:text-6xl font-black tracking-tighter text-white mb-4 font-mono">
            Your digital<br />
            <span className="text-emerald-400">piggy bank.</span>
          </h1>

          <p className="text-zinc-400 text-lg max-w-md mx-auto mb-10 leading-relaxed">
            Save little by little, withdraw when you need. Just like the piggy bank on your shelf — but smarter.
          </p>

          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <Link
              href="/register"
              className="w-full sm:w-auto bg-emerald-500 hover:bg-emerald-400 text-zinc-950 font-bold px-8 py-3.5 rounded-xl transition text-sm tracking-wide"
            >
              Start saving today
            </Link>
            <Link
              href="/login"
              className="w-full sm:w-auto border border-zinc-700 hover:border-zinc-500 text-zinc-300 hover:text-white font-medium px-8 py-3.5 rounded-xl transition text-sm"
            >
              I already have an account
            </Link>
          </div>
        </div>
      </main>

      {/* Feature strip */}
      <section className="px-6 pb-16 pt-4">
        <div className="max-w-2xl mx-auto grid grid-cols-1 sm:grid-cols-3 gap-4">
          {[
            { icon: "💰", title: "Save anytime", desc: "Drop money in whenever you want" },
            { icon: "📤", title: "Withdraw easily", desc: "Take out what you need, when you need it" },
            { icon: "📊", title: "Track it all", desc: "See every transaction at a glance" },
          ].map((f) => (
            <div
              key={f.title}
				className="bg-slate-900/70 border border-slate-700 rounded-xl p-5 text-center 
				shadow-lg shadow-black/30 backdrop-blur-md 
				hover:-translate-y-1 hover:shadow-xl hover:shadow-black/40 
				transition-all duration-300"
            >
              <div className="text-2xl mb-3">{f.icon}</div>
              <p className="text-white font-semibold text-sm mb-1">{f.title}</p>
              <p className="text-zinc-500 text-xs leading-relaxed">{f.desc}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Footer */}
      <footer className="text-center pb-8 text-zinc-600 text-xs">
        © {new Date().getFullYear()} Piggy. Built with care.
      </footer>
    </div>
  );
}