import React from "react";
// 1. Ganti impor dari lucide-react ke react-icons
// Kita akan menggunakan ikon dari set Font Awesome 6 (fa6) untuk mendapatkan logo X terbaru
import { FaXTwitter, FaLinkedinIn, FaInstagram, FaFacebookF } from "react-icons/fa6";

export default function Footer() {
    const menuItems = {
        Vintage: ["Homepage", "Technology", "Ataraxis Breast", "Resources & News"],
        Discover: ["Careers", "Blog", "News", "Events"],
        Help: ["FAQ", "Support", "Contact Us", "Portal"],
        Community: ["Forum", "Groups", "Ambassadors", "Partners"],
    };

    // 2. Gunakan komponen ikon dari react-icons
    const socialIcons = [
        { Icon: FaXTwitter, label: "Twitter", link: "#" },
        { Icon: FaLinkedinIn, label: "LinkedIn", link: "#" },
        { Icon: FaInstagram, label: "Instagram", link: "#" },
        { Icon: FaFacebookF, label: "Facebook", link: "#" },
    ];

    return (
        <footer className="bg-[#0f3a37] text-white py-16 px-6">
            <div className="max-w-7xl mx-auto flex flex-col md:flex-row md:justify-between gap-12 md:gap-0">
                {/* Logo & Social */}
                <div className="md:flex-1 max-w-sm">
                    <div className="flex items-center gap-3 mb-6">
                        <img
                            src="https://storage.googleapis.com/a1aa/image/3b6842ab-6b81-4b05-c444-1c40ae25855d.jpg"
                            alt="Logo"
                            className="w-6 h-6"
                        />
                        <span className="text-lg font-semibold tracking-wide">VINTAGE</span>
                    </div>
                    <p className="text-sm leading-relaxed mb-8 max-w-[280px]">
                        Empowering users with a unique vintage shopping experience.
                    </p>
                    <div className="flex space-x-6">
                        {/* Tidak ada perubahan yang diperlukan di sini */}
                        {socialIcons.map(({ Icon, label, link }) => (
                            <a
                                key={label}
                                aria-label={label}
                                href={link}
                                className="hover:text-gray-300"
                            >
                                <Icon size={20} />
                            </a>
                        ))}
                    </div>
                </div>

                {/* Menu */}
                <div className="flex flex-wrap gap-12 md:gap-20">
                    {Object.entries(menuItems).map(([category, items]) => (
                        <div key={category}>
                            <h3 className="text-sm font-semibold mb-4">{category}</h3>
                            <ul className="space-y-2 text-sm">
                                {items.map((item) => (
                                    <li key={item}>
                                        <a href="#" className="hover:text-gray-300">
                                            {item}
                                        </a>
                                    </li>
                                ))}
                            </ul>
                        </div>
                    ))}
                </div>
            </div>

            {/* Bottom bar */}
            <div className="mt-12 border-t border-white/20 pt-6 flex flex-col md:flex-row justify-between items-center text-sm text-gray-400">
                <div className="flex space-x-4">
                    <a href="#" className="hover:text-white">Privacy Policy</a>
                    <a href="#" className="hover:text-white">Terms of Service</a>
                </div>
                <span>&copy; {new Date().getFullYear()} VINTAGE. All Rights Reserved.</span>
            </div>
        </footer>
    );
}