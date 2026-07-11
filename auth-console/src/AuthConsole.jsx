import React, { useState } from "react";

const FONT_LINK_ID = "auth-console-fonts";

function useFonts() {
  React.useEffect(() => {
    if (document.getElementById(FONT_LINK_ID)) return;
    const link = document.createElement("link");
    link.id = FONT_LINK_ID;
    link.rel = "stylesheet";
    link.href =
      "https://fonts.googleapis.com/css2?family=Fraunces:opsz,wght@9..144,400;9..144,600;9..144,700&family=IBM+Plex+Mono:wght@400;500;600&family=IBM+Plex+Sans:wght@400;500;600&display=swap";
    document.head.appendChild(link);
  }, []);
}

const TABS = [
  { id: "register", label: "Register" },
  { id: "login", label: "Login" },
  { id: "profile", label: "View Profile" },
  { id: "complete", label: "Complete Profile" },
];

export default function AuthConsole() {
  useFonts();

  const [apiBase, setApiBase] = useState("http://localhost:8080");
  const [activeTab, setActiveTab] = useState("register");
  const [log, setLog] = useState([]);
  const [token, setToken] = useState(null);

  const [reg, setReg] = useState({ name: "", email: "", password: "" });
  const [regMsg, setRegMsg] = useState(null);
  const [regBusy, setRegBusy] = useState(false);

  const [login, setLogin] = useState({ email: "", password: "" });
  const [loginMsg, setLoginMsg] = useState(null);
  const [loginBusy, setLoginBusy] = useState(false);
  const [badge, setBadge] = useState(null);

  const [profile, setProfile] = useState(null);
  const [profileMsg, setProfileMsg] = useState(null);
  const [profileBusy, setProfileBusy] = useState(false);

  const [cp, setCp] = useState({ age: "", gender: "", leetcode: "" });
  const [completeMsg, setCompleteMsg] = useState(null);
  const [completeBusy, setCompleteBusy] = useState(false);

  function pushLog(method, path, ok, note) {
    setLog((prev) => [
      { time: new Date().toLocaleTimeString(), method, path, ok, note },
      ...prev,
    ]);
  }

  async function callApi(method, path, body, useAuth) {
    const headers = { "Content-Type": "application/json" };
    if (useAuth && token) headers["Authorization"] = "Bearer " + token;
    const opts = { method, headers };
    if (body !== undefined) opts.body = JSON.stringify(body);
    const res = await fetch(apiBase.replace(/\/$/, "") + path, opts);
    let data = null;
    try {
      data = await res.json();
    } catch (e) {
      /* no body */
    }
    return { ok: res.ok, status: res.status, data };
  }

  async function submitRegister() {
    if (!reg.name.trim() || !reg.email.trim() || !reg.password) {
      setRegMsg({ ok: false, text: "Name, email and password are all required." });
      return;
    }
    setRegBusy(true);
    try {
      const { ok, status, data } = await callApi("POST", "/register", {
        name: reg.name.trim(),
        email: reg.email.trim(),
        password: reg.password,
      });
      if (ok) {
        setRegMsg({
          ok: true,
          text: typeof data === "string" ? data : "Registered successfully. You can log in now.",
        });
        pushLog("POST", "/register", true, "HTTP " + status);
      } else {
        const errText = data && data.error ? data.error : "HTTP " + status;
        setRegMsg({ ok: false, text: "Registration failed: " + errText });
        pushLog("POST", "/register", false, errText);
      }
    } catch (e) {
      setRegMsg({ ok: false, text: `Could not reach the service at ${apiBase}. Is it running?` });
      pushLog("POST", "/register", false, "network error");
    }
    setRegBusy(false);
  }

  async function submitLogin() {
    if (!login.email.trim() || !login.password) {
      setLoginMsg({ ok: false, text: "Email and password are required." });
      return;
    }
    setLoginBusy(true);
    try {
      const { ok, status, data } = await callApi("POST", "/login", {
        email: login.email.trim(),
        password: login.password,
      });
      if (ok && data && data.token) {
        setToken(data.token);
        setLoginMsg({ ok: true, text: "Access granted." });
        pushLog("POST", "/login", true, "HTTP " + status);
        setBadge({ email: login.email.trim(), token: data.token });
      } else {
        const errText = data && data.error ? data.error : "HTTP " + status;
        setLoginMsg({ ok: false, text: "Login failed: " + errText });
        pushLog("POST", "/login", false, errText);
        setBadge(null);
      }
    } catch (e) {
      setLoginMsg({ ok: false, text: `Could not reach the service at ${apiBase}. Is it running?` });
      pushLog("POST", "/login", false, "network error");
    }
    setLoginBusy(false);
  }

  async function fetchProfile() {
    if (!token) {
      setProfileMsg({ ok: false, text: "Log in first — this route requires a token." });
      return;
    }
    setProfileBusy(true);
    try {
      const { ok, status, data } = await callApi("GET", "/profile", undefined, true);
      if (ok && data) {
        setProfileMsg({ ok: true, text: "Profile loaded." });
        pushLog("GET", "/profile", true, "HTTP " + status);
        setProfile(data);
      } else {
        setProfileMsg({ ok: false, text: `Could not load profile (HTTP ${status}).` });
        pushLog("GET", "/profile", false, "HTTP " + status);
      }
    } catch (e) {
      setProfileMsg({ ok: false, text: `Could not reach the service at ${apiBase}.` });
      pushLog("GET", "/profile", false, "network error");
    }
    setProfileBusy(false);
  }

  async function submitComplete() {
    if (!token) {
      setCompleteMsg({ ok: false, text: "Log in first — this route requires a token." });
      return;
    }
    const body = {};
    if (cp.age !== "") body.age = { Int64: parseInt(cp.age, 10), Valid: true };
    if (cp.gender.trim() !== "") body.gender = { String: cp.gender.trim(), Valid: true };
    if (cp.leetcode.trim() !== "") body.leetcode = { String: cp.leetcode.trim(), Valid: true };

    setCompleteBusy(true);
    try {
      const { ok, status, data } = await callApi("POST", "/complete_profile", body, true);
      if (ok) {
        setCompleteMsg({ ok: true, text: typeof data === "string" ? data : "Profile updated." });
        pushLog("POST", "/complete_profile", true, "HTTP " + status);
      } else {
        const errText = data && data.error ? data.error : "HTTP " + status;
        setCompleteMsg({ ok: false, text: "Update failed: " + errText });
        pushLog("POST", "/complete_profile", false, errText);
      }
    } catch (e) {
      setCompleteMsg({ ok: false, text: `Could not reach the service at ${apiBase}.` });
      pushLog("POST", "/complete_profile", false, "network error");
    }
    setCompleteBusy(false);
  }

  return (
    <div style={styles.wrap}>
      <style>{css}</style>

      {/* LEFT PANEL */}
      <div style={styles.terminal}>
        <div>
          <h1 style={styles.terminalH1}>Access Console</h1>
          <div style={styles.terminalSub}>User Auth Service — Control Panel</div>
        </div>

        <div style={styles.apiConfig}>
          <label style={styles.terminalLabel}>Service address</label>
          <input
            style={styles.apiInput}
            type="text"
            spellCheck={false}
            value={apiBase}
            onChange={(e) => setApiBase(e.target.value)}
          />
        </div>

        <div style={styles.sessionBadge}>
          <div style={styles.terminalLabel}>Session</div>
          <div style={token ? styles.statusPillOn : styles.statusPill}>
            <span style={styles.dot} /> {token ? "Token active" : "No active token"}
          </div>
          {token && <div style={styles.tokenPreview}>{token}</div>}
        </div>

        <div style={styles.terminalLabel}>Request log</div>
        <div style={styles.log}>
          {log.length === 0 && <div style={styles.logEmpty}>Requests you make will appear here.</div>}
          {log.map((entry, i) => (
            <div
              key={i}
              style={{
                ...styles.logEntry,
                borderLeftColor: entry.ok ? "#5FCBA6" : "#C1443C",
              }}
            >
              <div style={styles.logTime}>{entry.time}</div>
              <div style={styles.logMsg}>
                {entry.method} {entry.path} — {entry.note}
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* RIGHT PANEL */}
      <div style={styles.main}>
        <div style={styles.tabs}>
          {TABS.map((t) => (
            <div
              key={t.id}
              onClick={() => setActiveTab(t.id)}
              style={activeTab === t.id ? styles.tabActive : styles.tab}
            >
              {t.label}
            </div>
          ))}
        </div>

        <div style={styles.panel}>
          {activeTab === "register" && (
            <>
              <h2 style={styles.panelH2}>Create an account</h2>
              <div style={styles.desc}>
                Sends <code>POST /register</code> with your name, email and password. The service
                hashes the password before it's stored.
              </div>
              <Field label="Full name">
                <input
                  style={styles.input}
                  type="text"
                  placeholder="Jordan Lee"
                  value={reg.name}
                  onChange={(e) => setReg({ ...reg, name: e.target.value })}
                />
              </Field>
              <Field label="Email">
                <input
                  style={styles.input}
                  type="email"
                  placeholder="jordan@example.com"
                  value={reg.email}
                  onChange={(e) => setReg({ ...reg, email: e.target.value })}
                />
              </Field>
              <Field label="Password">
                <input
                  style={styles.input}
                  type="password"
                  placeholder="At least 8 characters"
                  value={reg.password}
                  onChange={(e) => setReg({ ...reg, password: e.target.value })}
                />
              </Field>
              <button style={styles.btn} disabled={regBusy} onClick={submitRegister}>
                Register
              </button>
              <Msg msg={regMsg} />
            </>
          )}

          {activeTab === "login" && (
            <>
              <h2 style={styles.panelH2}>Log in</h2>
              <div style={styles.desc}>
                Sends <code>POST /login</code>. On success the service returns a signed JWT, issued
                as your access badge below.
              </div>
              <Field label="Email">
                <input
                  style={styles.input}
                  type="email"
                  placeholder="jordan@example.com"
                  value={login.email}
                  onChange={(e) => setLogin({ ...login, email: e.target.value })}
                />
              </Field>
              <Field label="Password">
                <input
                  style={styles.input}
                  type="password"
                  placeholder="Your password"
                  value={login.password}
                  onChange={(e) => setLogin({ ...login, password: e.target.value })}
                />
              </Field>
              <button style={styles.btn} disabled={loginBusy} onClick={submitLogin}>
                Log in
              </button>
              <Msg msg={loginMsg} />

              {badge && (
                <div style={styles.badgeWrap}>
                  <div style={styles.badge}>
                    <div style={styles.badgeMain}>
                      <div style={styles.badgeEyebrow}>Access granted</div>
                      <div style={styles.badgeName}>{badge.email}</div>
                      <div style={styles.badgeId}>token issued</div>
                      <div style={styles.badgeToken}>{badge.token}</div>
                    </div>
                    <div style={styles.badgeStub}>Bearer</div>
                  </div>
                </div>
              )}
            </>
          )}

          {activeTab === "profile" && (
            <>
              <h2 style={styles.panelH2}>View profile</h2>
              <div style={styles.desc}>
                Sends <code>GET /profile</code> with your token in the <code>Authorization</code>{" "}
                header. Requires logging in first.
              </div>
              <button style={styles.btn} disabled={profileBusy} onClick={fetchProfile}>
                Fetch profile
              </button>
              <Msg msg={profileMsg} />
              {profile && (
                <div style={styles.profileGrid}>
                  <ProfileItem k="ID" v={profile.id ?? "—"} />
                  <ProfileItem k="Name" v={profile.name || "—"} />
                  <ProfileItem k="Email" v={profile.email || "—"} />
                  <ProfileItem
                    k="Age"
                    v={profile.age && profile.age.Valid ? profile.age.Int64 : "—"}
                  />
                  <ProfileItem
                    k="Gender"
                    v={profile.gender && profile.gender.Valid ? profile.gender.String : "—"}
                  />
                  <ProfileItem
                    k="Leetcode"
                    v={profile.leetcode && profile.leetcode.Valid ? profile.leetcode.String : "—"}
                  />
                </div>
              )}
            </>
          )}

          {activeTab === "complete" && (
            <>
              <h2 style={styles.panelH2}>Complete profile</h2>
              <div style={styles.desc}>
                Sends <code>POST /complete_profile</code> with your token, filling in age, gender
                and Leetcode handle. Requires logging in first.
              </div>
              <div style={styles.row2}>
                <Field label="Age">
                  <input
                    style={styles.input}
                    type="number"
                    min="0"
                    placeholder="28"
                    value={cp.age}
                    onChange={(e) => setCp({ ...cp, age: e.target.value })}
                  />
                </Field>
                <Field label="Gender">
                  <input
                    style={styles.input}
                    type="text"
                    placeholder="e.g. woman, man, non-binary"
                    value={cp.gender}
                    onChange={(e) => setCp({ ...cp, gender: e.target.value })}
                  />
                </Field>
              </div>
              <Field label="Leetcode profile">
                <input
                  style={styles.input}
                  type="text"
                  placeholder="leetcode.com/u/jordan"
                  value={cp.leetcode}
                  onChange={(e) => setCp({ ...cp, leetcode: e.target.value })}
                />
              </Field>
              <button style={styles.btn} disabled={completeBusy} onClick={submitComplete}>
                Save profile
              </button>
              <Msg msg={completeMsg} />
            </>
          )}
        </div>
      </div>
    </div>
  );
}

function Field({ label, children }) {
  return (
    <div style={styles.field}>
      <label style={styles.fieldLabel}>{label}</label>
      {children}
    </div>
  );
}

function Msg({ msg }) {
  if (!msg) return null;
  return (
    <div style={msg.ok ? styles.msgOk : styles.msgErr}>{msg.text}</div>
  );
}

function ProfileItem({ k, v }) {
  return (
    <div style={styles.profileItem}>
      <div style={styles.profileK}>{k}</div>
      <div style={styles.profileV}>{v}</div>
    </div>
  );
}

const colors = {
  paper: "#F0EEE6",
  panel: "#E8E4D8",
  ink: "#1B2430",
  inkSoft: "#4C5666",
  line: "#C9C3B2",
  clearanceRed: "#C1443C",
  verifiedTeal: "#2F6F62",
  verifiedTealBg: "#DCE8E3",
};

const mono = "'IBM Plex Mono', monospace";
const sans = "'IBM Plex Sans', sans-serif";
const serif = "'Fraunces', serif";

const styles = {
  wrap: {
    display: "grid",
    gridTemplateColumns: "300px 1fr",
    minHeight: "100vh",
    background: colors.paper,
    color: colors.ink,
    fontFamily: sans,
  },
  terminal: {
    background: colors.ink,
    color: "#D7DCE3",
    padding: "28px 24px",
    display: "flex",
    flexDirection: "column",
    borderRight: `1px solid ${colors.line}`,
  },
  terminalH1: {
    fontFamily: serif,
    fontWeight: 700,
    fontSize: 22,
    letterSpacing: "0.02em",
    color: "#F0EEE6",
    margin: "0 0 4px 0",
  },
  terminalSub: {
    fontFamily: mono,
    fontSize: 11,
    color: "#8A94A3",
    textTransform: "uppercase",
    letterSpacing: "0.12em",
    marginBottom: 22,
  },
  apiConfig: {
    marginBottom: 22,
    paddingBottom: 18,
    borderBottom: "1px dashed #3A4456",
  },
  terminalLabel: {
    display: "block",
    fontFamily: mono,
    fontSize: 10,
    textTransform: "uppercase",
    letterSpacing: "0.1em",
    color: "#8A94A3",
    marginBottom: 8,
  },
  apiInput: {
    width: "100%",
    background: "#0E1520",
    border: "1px solid #3A4456",
    color: "#D7DCE3",
    fontFamily: mono,
    fontSize: 12,
    padding: "8px 10px",
    borderRadius: 3,
  },
  sessionBadge: {
    marginBottom: 22,
    paddingBottom: 18,
    borderBottom: "1px dashed #3A4456",
  },
  statusPill: {
    display: "inline-flex",
    alignItems: "center",
    gap: 6,
    fontFamily: mono,
    fontSize: 11,
    padding: "5px 10px",
    borderRadius: 20,
    background: "#2A2015",
    color: "#D9A441",
  },
  statusPillOn: {
    display: "inline-flex",
    alignItems: "center",
    gap: 6,
    fontFamily: mono,
    fontSize: 11,
    padding: "5px 10px",
    borderRadius: 20,
    background: "#152A22",
    color: "#5FCBA6",
  },
  dot: {
    width: 6,
    height: 6,
    borderRadius: "50%",
    background: "currentColor",
    display: "inline-block",
  },
  tokenPreview: {
    fontFamily: mono,
    fontSize: 10,
    color: "#5C6779",
    marginTop: 8,
    wordBreak: "break-all",
    lineHeight: 1.5,
  },
  log: {
    flex: 1,
    fontFamily: mono,
    fontSize: 11,
    lineHeight: 1.7,
    overflowY: "auto",
    maxHeight: 340,
  },
  logEmpty: { color: "#5C6779", fontStyle: "italic" },
  logEntry: {
    marginBottom: 10,
    paddingLeft: 14,
    borderLeft: "2px solid #3A4456",
  },
  logTime: { color: "#5C6779", fontSize: 10 },
  logMsg: { color: "#D7DCE3" },
  main: { padding: "40px 48px", maxWidth: 640 },
  tabs: { display: "flex", gap: 6, marginBottom: 30, flexWrap: "wrap" },
  tab: {
    fontFamily: mono,
    fontSize: 11,
    letterSpacing: "0.06em",
    textTransform: "uppercase",
    padding: "9px 16px 8px",
    background: colors.panel,
    border: `1px solid ${colors.line}`,
    borderBottom: "none",
    borderRadius: "6px 6px 0 0",
    color: colors.inkSoft,
    cursor: "pointer",
    position: "relative",
    top: 1,
  },
  tabActive: {
    fontFamily: mono,
    fontSize: 11,
    letterSpacing: "0.06em",
    textTransform: "uppercase",
    padding: "9px 16px 8px",
    background: colors.paper,
    border: `1px solid ${colors.line}`,
    borderBottom: `1px solid ${colors.paper}`,
    borderRadius: "6px 6px 0 0",
    color: colors.ink,
    fontWeight: 600,
    cursor: "pointer",
    position: "relative",
    top: 1,
  },
  panel: {
    border: `1px solid ${colors.line}`,
    borderRadius: "0 8px 8px 8px",
    background: "#fff",
    padding: "30px 32px",
  },
  panelH2: { fontFamily: serif, fontSize: 24, fontWeight: 600, margin: "0 0 6px 0" },
  desc: { color: colors.inkSoft, fontSize: 13.5, marginBottom: 24, lineHeight: 1.5 },
  field: { marginBottom: 16 },
  fieldLabel: { display: "block", fontSize: 12, fontWeight: 600, color: colors.inkSoft, marginBottom: 6 },
  input: {
    width: "100%",
    padding: "10px 12px",
    border: `1px solid ${colors.line}`,
    borderRadius: 5,
    fontFamily: sans,
    fontSize: 14,
    background: "#FBFAF6",
    color: colors.ink,
  },
  row2: { display: "grid", gridTemplateColumns: "1fr 1fr", gap: 14 },
  btn: {
    fontFamily: mono,
    fontSize: 12,
    letterSpacing: "0.06em",
    textTransform: "uppercase",
    fontWeight: 600,
    padding: "11px 20px",
    borderRadius: 5,
    border: "none",
    cursor: "pointer",
    background: colors.ink,
    color: colors.paper,
  },
  msgOk: {
    marginTop: 16,
    padding: "11px 14px",
    borderRadius: 5,
    fontSize: 13,
    fontFamily: mono,
    background: colors.verifiedTealBg,
    color: colors.verifiedTeal,
    border: "1px solid #B9D3CB",
  },
  msgErr: {
    marginTop: 16,
    padding: "11px 14px",
    borderRadius: 5,
    fontSize: 13,
    fontFamily: mono,
    background: "#F6DEDB",
    color: colors.clearanceRed,
    border: "1px solid #E5B7B1",
  },
  badgeWrap: { marginTop: 26 },
  badge: {
    display: "flex",
    border: `1.5px solid ${colors.ink}`,
    borderRadius: 10,
    overflow: "hidden",
    background: "#fff",
    position: "relative",
  },
  badgeMain: { flex: 1, padding: "20px 22px", position: "relative" },
  badgeEyebrow: {
    fontFamily: mono,
    fontSize: 10,
    letterSpacing: "0.12em",
    textTransform: "uppercase",
    color: colors.verifiedTeal,
    marginBottom: 6,
  },
  badgeName: { fontFamily: serif, fontSize: 20, fontWeight: 600, marginBottom: 2 },
  badgeId: { fontFamily: mono, fontSize: 12, color: colors.inkSoft },
  badgeToken: {
    marginTop: 12,
    fontFamily: mono,
    fontSize: 9.5,
    color: "#8A94A3",
    wordBreak: "break-all",
    lineHeight: 1.5,
    maxHeight: 48,
    overflow: "hidden",
  },
  badgeStub: {
    width: 96,
    background: colors.ink,
    color: colors.paper,
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    writingMode: "vertical-rl",
    fontFamily: mono,
    fontSize: 11,
    letterSpacing: "0.15em",
    textTransform: "uppercase",
  },
  profileGrid: { marginTop: 22, display: "grid", gridTemplateColumns: "1fr 1fr", gap: 14 },
  profileItem: { padding: "12px 14px", background: colors.panel, borderRadius: 6 },
  profileK: {
    fontFamily: mono,
    fontSize: 10,
    textTransform: "uppercase",
    letterSpacing: "0.08em",
    color: colors.inkSoft,
    marginBottom: 4,
  },
  profileV: { fontSize: 14, fontWeight: 500 },
};

const css = `
  input:focus { outline: none; box-shadow: 0 0 0 3px ${colors.verifiedTealBg}; border-color: ${colors.verifiedTeal} !important; }
  button:hover { opacity: 0.88; }
  button:active { transform: translateY(1px); }
  button:disabled { opacity: 0.5; cursor: not-allowed; }
  @media (max-width: 860px) {
    div[style*="grid-template-columns: 300px"] { grid-template-columns: 1fr !important; }
  }
`;
