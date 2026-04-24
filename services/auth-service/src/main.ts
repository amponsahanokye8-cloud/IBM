import crypto from "crypto";
import express from "express";

type OtpRecord = {
  phone: string;
  code: string;
  expiresAt: number;
};

const app = express();
app.use(express.json());

const otpStore = new Map<string, OtpRecord>();
const otpTtlSec = Number(process.env.OTP_TTL_SEC ?? 300);

function makeOtpCode(): string {
  return String(Math.floor(100000 + Math.random() * 900000));
}

app.get("/health", (_req, res) => {
  res.json({ service: "auth-service", status: "ok" });
});

app.post("/v1/auth/request-otp", (req, res) => {
  const { phone } = req.body ?? {};
  if (typeof phone !== "string" || !phone.startsWith("+") || phone.length < 10) {
    return res.status(400).json({ error: "valid phone is required in E.164 format" });
  }

  const requestId = `otp_${crypto.randomUUID()}`;
  const code = makeOtpCode();
  const expiresAt = Date.now() + otpTtlSec * 1000;

  otpStore.set(requestId, { phone, code, expiresAt });

  return res.status(202).json({
    request_id: requestId,
    expires_in_sec: otpTtlSec,
    dev_code: process.env.NODE_ENV === "production" ? undefined : code
  });
});

app.post("/v1/auth/verify-otp", (req, res) => {
  const { request_id: requestId, otp_code: otpCode } = req.body ?? {};
  if (typeof requestId !== "string" || typeof otpCode !== "string") {
    return res.status(400).json({ error: "request_id and otp_code are required" });
  }

  const record = otpStore.get(requestId);
  if (!record) {
    return res.status(404).json({ error: "otp request not found" });
  }

  if (Date.now() > record.expiresAt) {
    otpStore.delete(requestId);
    return res.status(410).json({ error: "otp expired" });
  }

  if (record.code !== otpCode) {
    return res.status(401).json({ error: "invalid otp" });
  }

  otpStore.delete(requestId);

  return res.status(200).json({
    access_token: crypto.randomUUID(),
    refresh_token: crypto.randomUUID(),
    user: {
      id: crypto.randomUUID(),
      phone: record.phone
    }
  });
});

const port = Number(process.env.PORT ?? 3001);
app.listen(port, () => {
  console.log(`auth-service listening on :${port}`);
});
