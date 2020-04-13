module.exports = (req, res, next) => {
  if (!res.rawData && (typeof (res.rawData) !== "boolean") && (res.rawData !== 0)) return next();
  let rawData = {
    api: req.originalUrl,
    code: 200,
    message: "SUCCESS",
    method: req.method,
    version: "0.0.1",
    production: ISPROD,
    data: res.rawData,
    total: res.total || 0
  };

  res.status(200).send(rawData);
};