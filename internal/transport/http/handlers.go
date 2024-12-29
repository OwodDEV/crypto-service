package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// @Description Get USDT balance
// @Tags wallet
// @HeaderParam X-Request-ID string false "Optional request ID for tracing"
// @Param address path string true "Wallet Address" example(<br>ERC20 USDT: "0xe983fD1798689eee00c0Fb77e79B8f372DF41060", <br>TRC20 USDT: "TLSrrT5DiF5TkWPffJVQNwKE7SrctRCcpD")
// @Success 200 {object} models.GetWalletResp
// @Failure 400
// @Failure 500
// @Router /api/wallet/{address} [get]
func (s *Server) GetWalletHandler(c *fiber.Ctx) (err error) {
	ctx := c.UserContext()
	logger := slog.With(
		slog.String("request_id", ctx.Value("request_id").(string)),
	)

	address := c.Params("address")
	if address == "undefined" {
		err = errors.New("wallet address is empty")
		logger.Warn(err.Error())
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	resp, err := s.Service.GetWallet(ctx, address)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	c.JSON(resp)
	c.Status(http.StatusOK)
	return
}

// @Description Get USDT transaction details
// @Tags transaction
// @HeaderParam X-Request-ID string false "Optional request ID for tracing"
// @Param hash path string true "Transaction Hash" example(<br>ERC20 USDT: "0xec1d31abdcb80d24d0d823b35f93ed30c837d26364928e3b1b97b3c1cdd7fe69", <br>TRC20 USDT: "d6d1cc1ab403bc0febfb69d7be0bd8bd2fc03e2a03c4e2bdfd74560bd66109be")
// @Success 200 {object} models.GetTransactionResp
// @Failure 400
// @Failure 500
// @Router /api/transaction/{hash} [get]
func (s *Server) GetTransactionHandler(c *fiber.Ctx) (err error) {
	ctx := c.UserContext()
	logger := slog.With(
		slog.String("request_id", ctx.Value("request_id").(string)),
	)

	hash := c.Params("hash")
	if hash == "undefined" {
		err = errors.New("transaction hash is required")
		logger.Warn(err.Error())
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	resp, err := s.Service.GetTransaction(ctx, hash)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	c.JSON(resp)
	c.Status(http.StatusOK)
	return
}
