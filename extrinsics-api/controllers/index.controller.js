const { successResponse } = require("../helpers/methods")
const { ApiPromise, WsProvider, Keyring } = require("@polkadot/api")

/**
 *
 * @param req
 * @param res
 * @returns {Promise<void>}
 */
exports.index = async (req, res) => {
    const wsProvider = new WsProvider("ws://127.0.0.1:9944")
    const api = await ApiPromise.create({ provider: wsProvider })
    const keyring = new Keyring({ type: "sr25519" })

    const Alice = keyring.addFromUri("//Alice")
    const d = await api.tx.balances
        .transfer("5FNssn3gBGoSuJWPEgn1tHAizh711sG6VjmRtEwmS6w8XQab", 100000000000000)
        .signAndSend(Alice)
    res.send(
        successResponse("Express JS API Boiler Plate working like a charm...", {
            data: d
        })
    )
}

/**
 *
 * @param req
 * @param res
 * @returns {Promise<void>}
 */
exports.indexPost = async (req, res) => {
    res.send(
        successResponse("Express JS API Boiler Plate post api working like a charm...", {
            data: "here comes you payload...",
            request: req.body
        })
    )
}
