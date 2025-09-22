import React from "react";

export default function Navbar() {
    return (
        <header className="bg-white sticky top-0 z-50 ">
            <nav className="w-full  border-b border-gray-200 px-4 py-2 flex items-center justify-between max-w-7xl mx-auto">
                {/* Logo */}
                <div className="flex items-center space-x-2">
                    <img
                        src="https://storage.googleapis.com/a1aa/image/5a06ef6b-27fd-470a-4532-63fe4fa2df8d.jpg"
                        alt="Vintage logo"
                        className="w-6 h-6"
                    />
                    <span className="text-teal-700 font-semibold text-lg select-none">
                        Vintage
                    </span>
                </div>

                {/* Search */}
                <div className="flex-1 mx-4 max-w-xl">
                    <input
                        type="search"
                        placeholder="Search for items"
                        className="w-full border border-gray-300 rounded-lg py-2 px-3 text-gray-700 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-teal-600 focus:border-transparent transition"
                    />
                </div>

                {/* Actions */}
                <div className="flex items-center space-x-4 whitespace-nowrap">
                    <button
                        type="button"
                        className="text-teal-600 font-medium text-sm px-4 py-1 rounded border border-teal-600 hover:bg-teal-50 focus:outline-none focus:ring-2 focus:ring-teal-600 transition"
                    >
                        Login
                    </button>
                    <button
                        type="button"
                        className="bg-teal-700 text-white font-medium text-sm px-4 py-1 rounded hover:bg-teal-800 focus:outline-none focus:ring-2 focus:ring-teal-700 transition"
                    >
                        Sign up
                    </button>
                    <div className="text-gray-700 font-medium text-sm cursor-pointer select-none flex items-center space-x-1">
                        <span>EN</span>
                        <i className="fas fa-chevron-down text-xs"></i>
                    </div>
                </div>
            </nav>
        </header>
    );
}
