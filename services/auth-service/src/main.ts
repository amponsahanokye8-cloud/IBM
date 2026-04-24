import express from "express";

const app = express();
app.use(express.json());

app.get("/health", (_req, res) => {
  res.json({ service: "auth-service", status: "ok" });
});

app.post("/v1/auth/request-otp", (req, res) => {
  const { phone } = req.body ?? {};
  if (!phone) {
    return res.status(400).json({ error: "phone is required" });
  }

  return res.status(202).json({
    request_id: "otp_req_placeholder",
    expires_in_sec: 300
  });
});

const port = Number(process.env.PORT ?? 3001);
app.listen(port, () => {
  console.log(`auth-service listening on :${port}`);
});
