const methodNotAllowedHandler = async (req, res, next) => {
    res.status(405).send("Method not allowed");
}

module.exports = { methodNotAllowedHandler }