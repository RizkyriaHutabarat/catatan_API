package controller

import (
	"errors"
	"fmt"
	"net/http"

	inimodel "github.com/RizkyriaHutabarat/be_tb/Model"
	inimodul "github.com/RizkyriaHutabarat/be_tb/Module"

	"github.com/RizkyriaHutabarat/catatan_API/config"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Home(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"github_repo": "https://github.com/RizkyriaHutabarat/catatan_API",
		"message":     "You are at the root endpoint 😉",
		"success":     true,
	})
}

func GetAll(c *fiber.Ctx) error {
	ps := inimodul.GetAllCatatan(config.Ulbimongoconn, "catatan")
	return c.JSON(fiber.Map{
		"status": http.StatusOK,
		"data":   ps,
	})
}

func GetCatatanID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Wrong parameter",
		})
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Invalid id parameter",
		})
	}

	ps, err := inimodul.GetCatatanFromID(objID, config.Ulbimongoconn, "catatan")
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"status":  http.StatusNotFound,
				"message": fmt.Sprintf("No data found for id %s", id),
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": fmt.Sprintf("Error retrieving data for id %s", id),
		})
	}
	return c.JSON(ps)
}

func InsertData(c *fiber.Ctx) error {
	db := config.Ulbimongoconn
	var catatan inimodel.Catatan
	if err := c.BodyParser(&catatan); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}
	insertedID, err := inimodul.InsertCatatan(db, "catatan",
		catatan.Judul_Tugas,
		catatan.Matkul,
		catatan.Deskripsi_Tugas,
		catatan.Tanggal_Deadline,
		catatan.Tanggal_Submit,
	)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":      http.StatusOK,
		"message":     "Data berhasil disimpan.",
		"inserted_id": insertedID,
	})
}

func UpdateData(c *fiber.Ctx) error {
	db := config.Ulbimongoconn
	id := c.Params("id")
	var catatan inimodel.Catatan
	if err := c.BodyParser(&catatan); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}
	err = inimodul.UpdateCatatan(db, "catatan", oid,
		catatan.Judul_Tugas,
		catatan.Matkul,
		catatan.Deskripsi_Tugas,
		catatan.Tanggal_Deadline,
		catatan.Tanggal_Submit,
	)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Data berhasil diupdate.",
	})
}

func DeleteCatatan(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Wrong parameter",
		})
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Invalid id parameter",
		})
	}

	err = inimodul.DeleteCatatanByID(objID, config.Ulbimongoconn, "catatan")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": fmt.Sprintf("Error deleting data for id %s", id),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": fmt.Sprintf("Data with id %s deleted successfully", id),
	})
}

func Login(c *fiber.Ctx) error {
	db := config.Ulbimongoconn
	var user inimodel.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	authenticated, token, err := inimodul.Login(user.Username, user.Password, db, "user")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	if authenticated {
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"status":  http.StatusOK,
			"message": "Login successful",
			"token":   token,
		})
	}

	return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
		"status":  http.StatusUnauthorized,
		"message": "Invalid credentials",
	})
}