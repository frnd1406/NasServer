import { useState, useEffect } from "react";
import PropTypes from "prop-types";
import { Link, Outlet, useLocation, useNavigate } from "react-router-dom";
import { clearAuth, getAuth, isAuthenticated } from "./utils/auth";
import {
  LayoutDashboard,
  FolderOpen,
  Database,
  Settings,
  LogOut,
  Search,
  Bell,
  CloudLightning,
  Menu,
  X,
  Activity,
  MessageSquare,
  ChevronRight,
  Sparkles
} from "lucide-react";

// SidebarItem Component for reusable navigation items
const SidebarItem = ({ icon: Icon, label, active, onClick }) => (
  <button
    onClick={onClick}
    className={`flex items-center gap-3 w-full p-3 rounded-xl transition-all duration-300 group
      ${active
        ? 'bg-blue-500/20 text-blue-400 border border-blue-500/30 shadow-[0_0_15px_rgba(59,130,246,0.3)]'
        : 'text-slate-400 hover:bg-white/5 hover:text-white hover:pl-4'
      }`}
  >
    <Icon size={20} strokeWidth={1.5} />
    <span className="font-medium text-sm tracking-wide">{label}</span>
    {active && <div className="ml-auto w-1 h-1 bg-blue-400 rounded-full shadow-[0_0_8px_rgba(59,130,246,0.8)]" />}
  </button>
);

SidebarItem.propTypes = {
  icon: PropTypes.elementType.isRequired,
  label: PropTypes.string.isRequired,
  active: PropTypes.bool.isRequired,
  onClick: PropTypes.func.isRequired,
};

export default function Layout({ title = "NAS AI v1.0.0" }) {
  const location = useLocation();
  const navigate = useNavigate();
  const { accessToken } = getAuth();
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [scrolled, setScrolled] = useState(false);

  // Scroll effect for header
  useEffect(() => {
    const handleScroll = () => setScrolled(window.scrollY > 20);
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  const handleLogout = () => {
    clearAuth();
    navigate("/login", { replace: true });
  };

  // Auth guard: redirect unauthenticated users to login
  useEffect(() => {
    if (!isAuthenticated() && location.pathname !== "/login") {
      navigate("/login", { replace: true, state: { from: location.pathname } });
    }
  }, [location.pathname, navigate]);

  const navLinks = [
    { to: "/dashboard", label: "Dashboard", icon: LayoutDashboard },
    { to: "/search", label: "Semantic Search", icon: Search },
    { to: "/metrics", label: "Metrics", icon: Activity },
    { to: "/files", label: "Files & Storage", icon: FolderOpen },
    { to: "/backups", label: "Backups", icon: Database },
  ];

  return (
    <div className="min-h-screen bg-[#0a0a0c] text-slate-200 font-sans selection:bg-blue-500/30 overflow-x-hidden relative">

      {/* Background Effects (Liquid/Glow) */}
      <div className="fixed inset-0 z-0 pointer-events-none overflow-hidden">
        {/* Blue Blob */}
        <div className="absolute top-[-10%] left-[-10%] w-[500px] h-[500px] bg-blue-600/20 rounded-full blur-[120px] animate-pulse-glow"></div>
        {/* Violet Blob */}
        <div className="absolute bottom-[-10%] right-[-5%] w-[600px] h-[600px] bg-violet-600/10 rounded-full blur-[130px]"></div>
        {/* Cyan Accent middle */}
        <div className="absolute top-[40%] left-[30%] w-[300px] h-[300px] bg-cyan-500/10 rounded-full blur-[100px] opacity-60"></div>
      </div>

      <div className="relative z-10 flex h-screen overflow-hidden">

        {/* Sidebar (Desktop & Mobile) */}
        <aside className={`
          fixed inset-y-0 left-0 z-50 w-72 transform transition-transform duration-300 ease-in-out
          lg:relative lg:translate-x-0
          ${mobileMenuOpen ? 'translate-x-0' : '-translate-x-full'}
          bg-[#0a0a0c]/80 backdrop-blur-2xl border-r border-white/5 flex flex-col
        `}>
          {/* Logo Area */}
          <div className="p-8 pb-4 flex items-center gap-3">
            <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-blue-500 to-violet-600 flex items-center justify-center shadow-lg shadow-blue-500/20">
              <CloudLightning size={22} className="text-white" />
            </div>
            <div>
              <h1 className="text-xl font-bold text-white tracking-wide">{title}</h1>
              <p className="text-[10px] text-blue-400 font-medium tracking-widest uppercase">System Online</p>
            </div>
          </div>

          {/* Navigation */}
          <nav className="flex-1 px-4 py-6 space-y-2 overflow-y-auto">
            <p className="px-4 text-xs font-semibold text-slate-500 uppercase tracking-wider mb-2">Main Menu</p>
            {navLinks.map((link) => {
              const active = location.pathname.startsWith(link.to);
              return (
                <Link key={link.to} to={link.to} className="block">
                  <SidebarItem
                    icon={link.icon}
                    label={link.label}
                    active={active}
                    onClick={() => setMobileMenuOpen(false)}
                  />
                </Link>
              );
            })}

            <p className="px-4 text-xs font-semibold text-slate-500 uppercase tracking-wider mt-8 mb-2">Preferences</p>
            <SidebarItem
              icon={Settings}
              label="Settings"
              active={location.pathname === '/settings'}
              onClick={() => {
                navigate('/settings');
                setMobileMenuOpen(false);
              }}
            />
          </nav>

          {/* User & Logout */}
          <div className="p-4 border-t border-white/5 bg-gradient-to-t from-black/40 to-transparent">
            {isAuthenticated() && (
              <>
                <div className="flex items-center gap-3 p-3 rounded-xl bg-white/5 border border-white/5 mb-3">
                  <div className="w-8 h-8 rounded-full bg-slate-700 overflow-hidden border border-white/10">
                    <div className="w-full h-full bg-gradient-to-tr from-slate-600 to-slate-500 flex items-center justify-center text-xs text-white font-bold">
                      {accessToken ? 'AU' : 'U'}
                    </div>
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-white truncate">User</p>
                    <p className="text-xs text-slate-400 truncate">Admin Access</p>
                  </div>
                </div>
                <button
                  onClick={handleLogout}
                  className="flex items-center gap-3 w-full p-3 rounded-xl text-rose-400 hover:bg-rose-500/10 transition-colors group"
                >
                  <LogOut size={20} className="group-hover:translate-x-1 transition-transform" />
                  <span className="font-medium text-sm">Logout</span>
                </button>
              </>
            )}
            {!accessToken && (
              <Link to="/login" className="flex items-center gap-3 w-full p-3 rounded-xl text-blue-400 hover:bg-blue-500/10 transition-colors">
                <span className="font-medium text-sm">Login</span>
              </Link>
            )}
          </div>
        </aside>

        {/* Main Content Area */}
        <main className="flex-1 h-full overflow-y-auto relative scroll-smooth">
          {/* Header */}
          <header className={`sticky top-0 z-40 px-6 py-4 flex items-center justify-between transition-all duration-300 ${scrolled ? 'bg-[#0a0a0c]/80 backdrop-blur-md border-b border-white/5' : ''}`}>
            <div className="flex items-center gap-4 lg:hidden">
              <button onClick={() => setMobileMenuOpen(!mobileMenuOpen)} className="p-2 text-slate-400 hover:text-white">
                {mobileMenuOpen ? <X size={24} /> : <Menu size={24} />}
              </button>
              <span className="text-lg font-bold text-white">{title}</span>
            </div>

            <div className="hidden lg:block">
              <h2 className="text-2xl font-semibold text-white">
                {location.pathname === '/dashboard' && 'Dashboard Overview'}
                {location.pathname === '/search' && 'Semantic Search'}
                {location.pathname === '/metrics' && 'System Metrics'}
                {location.pathname === '/files' && 'Files & Storage'}
                {location.pathname === '/backups' && 'Backups'}
                {location.pathname === '/settings' && 'Settings'}
              </h2>
              <p className="text-slate-400 text-sm">Welcome back to your neural hub.</p>
            </div>

            <div className="flex items-center gap-4">
              <div className="hidden md:flex items-center bg-white/5 border border-white/10 rounded-full px-4 py-2 w-64 focus-within:bg-white/10 focus-within:border-blue-500/50 transition-all group">
                <Search size={18} className="text-slate-400 group-focus-within:text-blue-400 transition-colors" />
                <input
                  type="text"
                  placeholder="Suche... (Enter)"
                  className="bg-transparent border-none outline-none text-sm text-white ml-2 w-full placeholder:text-slate-500"
                  onKeyDown={(e) => {
                    if (e.key === 'Enter' && e.target.value.trim()) {
                      navigate(`/search?q=${encodeURIComponent(e.target.value.trim())}`);
                      e.target.value = '';
                    }
                  }}
                />
              </div>
              <button className="relative p-2.5 rounded-full bg-white/5 border border-white/10 text-slate-300 hover:bg-white/10 transition-colors">
                <Bell size={20} />
                <span className="absolute top-2 right-2 w-2 h-2 bg-rose-500 rounded-full shadow-[0_0_8px_rgba(244,63,94,0.6)]"></span>
              </button>
            </div>
          </header>

          {/* Page Content */}
          <div className="p-6 lg:p-10 max-w-[1600px] mx-auto">
            <Outlet />
          </div>
        </main>
      </div>

      {/* Global Chat Widget */}
      <ChatWidget />
    </div>
  );
}

// ChatWidget Component - Floating AI Chat Overlay
function ChatWidget() {
  const [isOpen, setIsOpen] = useState(false);
  const [input, setInput] = useState('');
  const [messages, setMessages] = useState([
    { id: 1, role: 'ai', text: 'Hallo! Ich bin dein NAS.AI Assistant. Wie kann ich dir helfen?' }
  ]);

  const sendMessage = (e) => {
    e.preventDefault();
    if (!input.trim()) return;
    setMessages([...messages, { id: Date.now(), role: 'user', text: input }]);
    setInput('');
    // Simulate AI response
    setTimeout(() => {
      setMessages(prev => [...prev, {
        id: Date.now() + 1,
        role: 'ai',
        text: 'Ich analysiere deine Anfrage... Diese Funktion wird bald mit dem AI Knowledge Agent verbunden.'
      }]);
    }, 600);
  };

  return (
    <>
      {/* Floating Action Button */}
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="fixed bottom-6 right-6 p-4 bg-gradient-to-br from-blue-600 to-violet-600 text-white rounded-full shadow-[0_0_30px_rgba(79,70,229,0.3)] hover:shadow-[0_0_40px_rgba(79,70,229,0.5)] hover:scale-105 transition-all duration-300 z-50"
      >
        {isOpen ? <X size={24} /> : <MessageSquare size={24} />}
        {!isOpen && (
          <span className="absolute -top-1 -right-1 w-3 h-3 bg-emerald-500 rounded-full border-2 border-[#0a0a0c] animate-pulse" />
        )}
      </button>

      {/* Chat Panel */}
      {isOpen && (
        <div className="fixed bottom-24 right-6 w-96 h-[500px] bg-slate-900/95 backdrop-blur-xl border border-white/10 rounded-2xl shadow-2xl flex flex-col z-40 animate-in slide-in-from-bottom-10 fade-in duration-300">
          {/* Header */}
          <div className="p-4 border-b border-white/10 flex items-center bg-slate-800/30 rounded-t-2xl">
            <div className="w-2 h-2 rounded-full bg-emerald-500 mr-2 animate-pulse" />
            <span className="font-semibold text-white">NAS.AI Assistant</span>
            <Sparkles size={14} className="ml-2 text-violet-400" />
            <span className="ml-auto text-xs text-slate-500 bg-slate-800 px-2 py-1 rounded">Online</span>
          </div>

          {/* Messages */}
          <div className="flex-1 overflow-y-auto p-4 space-y-4">
            {messages.map((msg) => (
              <div key={msg.id} className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}>
                <div className={`max-w-[80%] rounded-2xl px-4 py-3 text-sm ${msg.role === 'user'
                  ? 'bg-blue-600 text-white rounded-br-none'
                  : 'bg-slate-800 border border-white/5 text-slate-200 rounded-bl-none'
                  }`}>
                  {msg.text}
                </div>
              </div>
            ))}
          </div>

          {/* Input */}
          <form onSubmit={sendMessage} className="p-4 border-t border-white/10 bg-slate-800/30 rounded-b-2xl">
            <div className="relative">
              <input
                type="text"
                className="w-full bg-slate-800 text-white pl-4 pr-10 py-3 rounded-xl border border-white/10 focus:outline-none focus:border-blue-500/50 transition-colors placeholder-slate-500 text-sm"
                placeholder="Frag etwas..."
                value={input}
                onChange={(e) => setInput(e.target.value)}
              />
              <button type="submit" className="absolute right-2 top-2 p-1.5 bg-blue-600 rounded-lg text-white hover:bg-blue-500 transition-colors">
                <ChevronRight size={16} />
              </button>
            </div>
          </form>
        </div>
      )}
    </>
  );
}

Layout.propTypes = {
  title: PropTypes.string,
};
