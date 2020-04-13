/*global ISPROD */

module.exports = (err, req, res, next) => {
  const resultData = {
    api: req.originalUrl,
    method: req.method,
    message: err && err.message ? err.message : err,
    version: "0.0.1",
    production: ISPROD,
    code: -200,
  };

  res.status(err.status || 200).send(resultData);
};
