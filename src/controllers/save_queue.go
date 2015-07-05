package controllers
import (
    "encoding/json"
	"utils"
	"log"
	"fmt"
	"os"
	"regexp"
	"encoding/pem"
	"errors"
	"time"
	"crypto/x509"
	"crypto/rsa"
	"crypto/rand"
	"crypto"
	"io/ioutil"
)

func (c *Controller) SaveQueue() (string, error) {

	var err error
	c.r.ParseForm()

	userId := []byte(c.r.FormValue("user_id"))
	if !utils.CheckInputData(userId, "int") {
		return `{"result":"incorrect userId"}`, nil
	}
	txTime := utils.StrToInt64(c.r.FormValue("time"))
	if !utils.CheckInputData(txTime, "int") {
		return `{"result":"incorrect time"}`, nil
	}
	txType_ := c.r.FormValue("type")
	if !utils.CheckInputData(txType_, "type") {
		return `{"result":"incorrect type"}`, nil
	}
	txType := utils.TypeInt(txType_)
	signature1 := c.r.FormValue("signature1")
	signature2 := c.r.FormValue("signature2")
	signature3 := c.r.FormValue("signature3")
	sign := utils.EncodeLengthPlusData([]byte(signature1))
	if len(signature2) > 0 {
		sign = append(sign, utils.EncodeLengthPlusData([]byte(signature2))...)
	}
	if len(signature3) > 0 {
		sign = append(sign, utils.EncodeLengthPlusData([]byte(signature3))...)
	}
	binSignatures := utils.EncodeLengthPlusData([]byte(sign))

	log.Println("txType_", txType_)

	var data []byte
	switch txType_ {
	case "new_user":
		publicKeyHex := c.r.FormValue("public_key")
		publicKey := utils.HexToBin([]byte(publicKeyHex))
		privateKey := c.r.FormValue("private_key")
		verifyData := map[string]string {publicKeyHex: "public_key", privateKey: "private_key"}
		err := CheckInputData(verifyData)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		if c.SessRestricted==0 {
			err = c.ExecSql(`
					INSERT INTO  `+c.MyPrefix+`my_new_users (
						public_key,
						private_key
					)
					VALUES (
						[hex],
						?
					)`, publicKeyHex, privateKey)
			if err != nil {
				return "", utils.ErrInfo(err)
			}
		}

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData([]byte(userId))...)
		data = append(data, utils.EncodeLengthPlusData([]byte(publicKey))...)
		data = append(data, binSignatures...)

	case "del_cf_project" :

		projectId := []byte(c.r.FormValue("project_id"));

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(projectId)...)
		data = append(data, binSignatures...)



	case "cf_comment" :

		projectId := []byte(c.r.FormValue("project_id"));
		langId := []byte(c.r.FormValue("lang_id"));
		comment := []byte(c.r.FormValue("comment"));

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(projectId)...)
		data = append(data, utils.EncodeLengthPlusData(langId)...)
		data = append(data, utils.EncodeLengthPlusData(comment)...)
		data = append(data, binSignatures...)




	case "new_credit" :

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData([]byte(c.r.FormValue("to_user_id")))...)
		data = append(data, utils.EncodeLengthPlusData([]byte(c.r.FormValue("amount")))...)
		data = append(data, utils.EncodeLengthPlusData([]byte(c.r.FormValue("currency_id")))...)
		data = append(data, utils.EncodeLengthPlusData([]byte(c.r.FormValue("pct")))...)
		data = append(data, binSignatures...)




	case "del_credit" :

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData([]byte(c.r.FormValue("credit_id")))...)
		data = append(data, binSignatures...)



	case "repayment_credit" :

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData([]byte(c.r.FormValue("credit_id")))...)
		data = append(data, utils.EncodeLengthPlusData([]byte(c.r.FormValue("amount")))...)
		data = append(data, binSignatures...)



	case "change_creditor" :

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData([]byte(c.r.FormValue("to_user_id")))...)
		data = append(data, utils.EncodeLengthPlusData([]byte(c.r.FormValue("credit_id")))...)
		data = append(data, binSignatures...)



	case "change_credit_part" :

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData([]byte(c.r.FormValue("pct")))...)
		data = append(data, binSignatures...)



	case "user_avatar" :

		name := []byte(c.r.FormValue("name"));
		avatar := []byte(c.r.FormValue("avatar"));

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(name)...)
		data = append(data, utils.EncodeLengthPlusData(avatar)...)
		data = append(data, binSignatures...)



	case "del_cf_funding" :

		fundingId := []byte(c.r.FormValue("funding_id"));

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(fundingId)...)
		data = append(data, binSignatures...)



	case "cf_project_change_category" :

		projectId := []byte(c.r.FormValue("project_id"));
		categoryId := []byte(c.r.FormValue("category_id"));

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(projectId)...)
		data = append(data, utils.EncodeLengthPlusData(categoryId)...)
		data = append(data, binSignatures...)



	case "new_cf_project" :

		currencyId := []byte(c.r.FormValue("currency_id"));
		amount := []byte(c.r.FormValue("amount"));
		endTime := []byte(c.r.FormValue("end_time"));
		latitude := []byte(c.r.FormValue("latitude"));
		longitude := []byte(c.r.FormValue("longitude"));
		categoryId := []byte(c.r.FormValue("category_id"));
		currencyName := []byte(c.r.FormValue("currency_name"));

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(currencyId)...)
		data = append(data, utils.EncodeLengthPlusData(amount)...)
		data = append(data, utils.EncodeLengthPlusData(endTime)...)
		data = append(data, utils.EncodeLengthPlusData(latitude)...)
		data = append(data, utils.EncodeLengthPlusData(longitude)...)
		data = append(data, utils.EncodeLengthPlusData(categoryId)...)
		data = append(data, utils.EncodeLengthPlusData(currencyName)...)
		data = append(data, binSignatures...)



	case "cf_project_data" :

		projectId := []byte(c.r.FormValue("project_id"));
		langId := []byte(c.r.FormValue("lang_id"));
		blurbImg := []byte(c.r.FormValue("blurb_img"));
		headImg := []byte(c.r.FormValue("head_img"));
		descriptionImg := []byte(c.r.FormValue("description_img"));
		picture := []byte(c.r.FormValue("picture"));
		videoType := []byte(c.r.FormValue("video_type"));
		videoUrlId := []byte(c.r.FormValue("video_url_id"));
		newsImg := []byte(c.r.FormValue("news_img"));
		links := []byte(c.r.FormValue("links"));
		hide := []byte(c.r.FormValue("hide"));

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(projectId)...)
		data = append(data, utils.EncodeLengthPlusData(langId)...)
		data = append(data, utils.EncodeLengthPlusData(blurbImg)...)
		data = append(data, utils.EncodeLengthPlusData(headImg)...)
		data = append(data, utils.EncodeLengthPlusData(descriptionImg)...)
		data = append(data, utils.EncodeLengthPlusData(picture)...)
		data = append(data, utils.EncodeLengthPlusData(videoType)...)
		data = append(data, utils.EncodeLengthPlusData(videoUrlId)...)
		data = append(data, utils.EncodeLengthPlusData(newsImg)...)
		data = append(data, utils.EncodeLengthPlusData(links)...)
		data = append(data, utils.EncodeLengthPlusData(hide)...)
		data = append(data, binSignatures...)



	case "new_miner" :

		race := []byte(c.r.FormValue("race"));
		country := []byte(c.r.FormValue("country"));
		latitude := []byte(c.r.FormValue("latitude"));
		longitude := []byte(c.r.FormValue("longitude"));
		host := []byte(c.r.FormValue("host"));
		faceHash := []byte(c.r.FormValue("face_hash"));
		profileHash := []byte(c.r.FormValue("profile_hash"));
		faceCoords := []byte(c.r.FormValue("face_coords"));
		profileCoords := []byte(c.r.FormValue("profile_coords"));
		videoType := []byte(c.r.FormValue("video_type"));
		videoUrlId := []byte(c.r.FormValue("video_url_id"));
		nodePublicKey := []byte(c.r.FormValue("node_public_key"));

		if len(race) == 0 || len(country) == 0 || len(latitude) == 0 || len(longitude) == 0 || len(host) == 0 || len(faceHash) == 0 || len(profileHash) == 0 || len(faceCoords) == 0 || len(profileCoords) == 0 || len(videoType) == 0 || len(videoUrlId) == 0 || len(nodePublicKey) == 0 {
			return "empty", nil
		}
		if string(videoType) == "null" || string(videoUrlId) == "null" {
			if _, err := os.Stat("public/"+utils.Int64ToStr(c.SessUserId)+"_user_video.mp4"); os.IsNotExist(err) {
				return "empty video", nil
			}
		}
		nodePublicKey = utils.HexToBin(nodePublicKey);
		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(race)...)
		data = append(data, utils.EncodeLengthPlusData(country)...)
		data = append(data, utils.EncodeLengthPlusData(latitude)...)
		data = append(data, utils.EncodeLengthPlusData(longitude)...)
		data = append(data, utils.EncodeLengthPlusData(host)...)
		data = append(data, utils.EncodeLengthPlusData(faceCoords)...)
		data = append(data, utils.EncodeLengthPlusData(profileCoords)...)
		data = append(data, utils.EncodeLengthPlusData(faceHash)...)
		data = append(data, utils.EncodeLengthPlusData(profileHash)...)
		data = append(data, utils.EncodeLengthPlusData(videoType)...)
		data = append(data, utils.EncodeLengthPlusData(videoUrlId)...)
		data = append(data, utils.EncodeLengthPlusData(nodePublicKey)...)
		data = append(data, binSignatures...)

		if c.SessRestricted == 0 {
			err := c.ExecSql(`UPDATE `+c.MyPrefix+`my_table
					SET node_voting_send_request = ?`, txTime)
			if err != nil {
				return "", utils.ErrInfo(err)
			}
		}


	case "votes_miner" : // голос за юзера, который хочет стать майнером

		voteId := []byte(c.r.FormValue("vote_id"));
		result := []byte(c.r.FormValue("result"));
		comment := []byte(c.r.FormValue("comment"));

		if c.SessRestricted == 0 {
			err := c.ExecSql(`INSERT INTO  `+c.MyPrefix+`my_tasks (
							type,
							id,
							time
						)
						VALUES (
							'miner',
							?,
							?
						)`, voteId, txTime)
			if err != nil {
				return "", utils.ErrInfo(err)
			}
		}

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(voteId)...)
		data = append(data, utils.EncodeLengthPlusData(result)...)
		data = append(data, utils.EncodeLengthPlusData(comment)...)
		data = append(data, binSignatures...)




	case "new_promised_amount" :

		currencyId := []byte(c.r.FormValue("currency_id"));
		amount := []byte(c.r.FormValue("amount"));
		videoType := []byte(c.r.FormValue("video_type"));
		videoUrlId := []byte(c.r.FormValue("video_url_id"));
		paymentSystemsIds := []byte(c.r.FormValue("payment_systems_ids"));

		verifyData := map[string]string {c.r.FormValue("currency_id"): "int", c.r.FormValue("amount"): "amount", c.r.FormValue("payment_systems_ids"): "payment_systems_ids"}
		err := CheckInputData(verifyData)
		if err != nil {
			return "", utils.ErrInfo(err)
		}

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(currencyId)...)
		data = append(data, utils.EncodeLengthPlusData(amount)...)
		data = append(data, utils.EncodeLengthPlusData(videoType)...)
		data = append(data, utils.EncodeLengthPlusData(videoUrlId)...)
		data = append(data, utils.EncodeLengthPlusData(paymentSystemsIds)...)
		data = append(data, binSignatures...)

		if c.SessRestricted == 0 {
			err = c.ExecSql(`
					INSERT INTO  `+c.MyPrefix+`my_promised_amount (
						currency_id,
						amount
					)
					VALUES (
						?,
						?
					)`, currencyId, amount)
			if err != nil {
				return "", utils.ErrInfo(err)
			}
		}



	case "change_promised_amount" :

		promisedAmountId := []byte(c.r.FormValue("promised_amount_id"));
		amount := []byte(c.r.FormValue("amount"));
		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(promisedAmountId)...)
		data = append(data, utils.EncodeLengthPlusData(amount)...)
		data = append(data, binSignatures...)



	case "mining" :

		promisedAmountId := []byte(c.r.FormValue("promised_amount_id"));
		amount := []byte(c.r.FormValue("amount"));
		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(promisedAmountId)...)
		data = append(data, utils.EncodeLengthPlusData(amount)...)
		data = append(data, binSignatures...)



	case "votes_promised_amount":

		promisedAmountId := []byte(c.r.FormValue("promised_amount_id"));
		result := []byte(c.r.FormValue("result"));
		comment := []byte(c.r.FormValue("comment"));

		if c.SessRestricted == 0 {
			err := c.ExecSql(`INSERT INTO  `+c.MyPrefix+`my_tasks (
							type,
							id,
							time
						)
						VALUES (
							'promised_amount',
							?,
							?
						)`, promisedAmountId, txTime)
			if err != nil {
				return "", utils.ErrInfo(err)
			}
		}

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(promisedAmountId)...)
		data = append(data, utils.EncodeLengthPlusData(result)...)
		data = append(data, utils.EncodeLengthPlusData(comment)...)
		data = append(data, binSignatures...)


	case "change_geolocation" :

		latitude := []byte(c.r.FormValue("latitude"));
		longitude := []byte(c.r.FormValue("longitude"));
		country := []byte(c.r.FormValue("country"));

		verifyData := map[string]string {c.r.FormValue("latitude"): "coordinate", c.r.FormValue("longitude"): "coordinate", c.r.FormValue("country"): "int"}
		err := CheckInputData(verifyData)
		if err != nil {
			return "", utils.ErrInfo(err)
		}

		if c.SessRestricted == 0 {
			err = c.ExecSql(`
				UPDATE `+c.MyPrefix+`my_table
				SET geolocation = ?,
					   location_country =  ?`, string(latitude)+", "+string(longitude), country)
			if err != nil {
				return "", utils.ErrInfo(err)
			}
		}

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(latitude)...)
		data = append(data, utils.EncodeLengthPlusData(longitude)...)
		data = append(data, utils.EncodeLengthPlusData(country)...)
		data = append(data, binSignatures...)



	case "del_promised_amount" :

		promisedAmountId := []byte(c.r.FormValue("promised_amount_id"));

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(promisedAmountId)...)
		data = append(data, binSignatures...)


	case "del_forex_order" :

		orderId := []byte(c.r.FormValue("order_id"));

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(orderId)...)
		data = append(data, binSignatures...)




	case "sendDc" :

		toUserId := []byte(c.r.FormValue("to_id"));
		currencyId := []byte(c.r.FormValue("currency_id"));
		amount := []byte(c.r.FormValue("amount"));
		commission := utils.StrToFloat64(c.r.FormValue("commission"));

		var arbitrators_ []string
		err := json.Unmarshal([]byte(c.r.PostFormValue("arbitrators")), &arbitrators_)
		if err != nil {
			return fmt.Sprintf("%q", err), err
		}
		arbitrators := make(map[int]string)
		for i:=0; i<len(arbitrators_); i++ {
			arbitrators[i] = arbitrators_[i]
		}

		var arbitrators_commissions_ []float64
		err = json.Unmarshal([]byte(c.r.PostFormValue("arbitrators_commissions")), &arbitrators_commissions_)
		if err != nil {
			return fmt.Sprintf("%q", err), err
		}
		arbitrators_commissions := make(map[int]float64)
		for i:=0; i<len(arbitrators_commissions_); i++ {
			arbitrators_commissions[i] = arbitrators_commissions_[i]
		}

		var arbitrators_commissions_sum float64
		for i:=0; i < 5; i++ {
			if len(arbitrators[i]) > 0 {
				if !utils.CheckInputData(arbitrators[i], "int") {
					return "incorrect arbitrators", nil
				}
				if ok, _ := regexp.MatchString(`^[0-9]{0,10}(\.[0-9]{0,2})?$`, utils.Float64ToStrPct(arbitrators_commissions[i])); !ok{
					return "incorrect arbitrator_commission", nil
				}
			} else {
				arbitrators[i] = "0"
				arbitrators_commissions[i] = 0
			}
			arbitrators_commissions_sum += arbitrators_commissions[i]
		}

		comment := []byte(c.r.FormValue("comment"));
		commentText := []byte(c.r.FormValue("comment_text"));

		verifyData := map[string]string {c.r.FormValue("to_id"): "int", c.r.FormValue("currency_id"): "int", c.r.FormValue("amount"): "amount", c.r.FormValue("commission"): "amount"}
		err = CheckInputData(verifyData)
		if err != nil {
			return "", utils.ErrInfo(err)
		}

		total_commission := commission + arbitrators_commissions_sum;
		if c.SessRestricted == 0 {
			// пишем транзакцкцию к себе в таблу
			err = c.ExecSql(`INSERT INTO
								`+c.MyPrefix+`my_dc_transactions (
									status,
									type,
									type_id,
									to_user_id,
									amount,
									commission,
									currency_id,
									comment,
									comment_status
								)
								VALUES (
									'pending',
									'from_user',
									?,
									?,
									?,
									?,
									?,
									?,
									'decrypted'
								)`, userId, toUserId, amount, total_commission, currencyId, commentText)
			if err != nil {
				return "", utils.ErrInfo(err)
			}
		}

		if len(comment) == 0 {
			comment = []byte("null");
		} else {
			comment = utils.HexToBin(comment);
		}

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(toUserId)...)
		data = append(data, utils.EncodeLengthPlusData(currencyId)...)
		data = append(data, utils.EncodeLengthPlusData(amount)...)
		data = append(data, utils.EncodeLengthPlusData([]byte(c.r.FormValue("commission")))...)
		data = append(data, utils.EncodeLengthPlusData([]byte(arbitrators[0]))...)
		data = append(data, utils.EncodeLengthPlusData([]byte(arbitrators[1]))...)
		data = append(data, utils.EncodeLengthPlusData([]byte(arbitrators[2]))...)
		data = append(data, utils.EncodeLengthPlusData([]byte(arbitrators[3]))...)
		data = append(data, utils.EncodeLengthPlusData([]byte(arbitrators[4]))...)
		data = append(data, utils.EncodeLengthPlusData(utils.Float64ToBytes(arbitrators_commissions[0]))...)
		data = append(data, utils.EncodeLengthPlusData(utils.Float64ToBytes(arbitrators_commissions[1]))...)
		data = append(data, utils.EncodeLengthPlusData(utils.Float64ToBytes(arbitrators_commissions[2]))...)
		data = append(data, utils.EncodeLengthPlusData(utils.Float64ToBytes(arbitrators_commissions[3]))...)
		data = append(data, utils.EncodeLengthPlusData(utils.Float64ToBytes(arbitrators_commissions[4]))...)
		data = append(data, utils.EncodeLengthPlusData(comment)...)
		data = append(data, binSignatures...)


	case "CfSendDc" :

		projectId := []byte(c.r.FormValue("to_id"));
		amount := []byte(c.r.FormValue("amount"));
		commission := []byte(c.r.FormValue("commission"));
		comment := []byte(c.r.FormValue("comment"));
		commentText := []byte(c.r.FormValue("comment_text"));

		verifyData := map[string]string {c.r.FormValue("to_id"): "int", c.r.FormValue("amount"): "amount", c.r.FormValue("commission"): "amount"}
		err := CheckInputData(verifyData)
		if err != nil {
			return "", utils.ErrInfo(err)
		}

		currencyId, err := c.Single(`SELECT currency_id
					FROM cf_projects
					WHERE id = ?`, projectId).Int64()
		if err != nil {
			return "", utils.ErrInfo(err)
		}

		if c.SessRestricted == 0 {
			// пишем транзакцкцию к сбе в таблу
			err = c.ExecSql(`INSERT INTO
							`+c.MyPrefix+`my_dc_transactions (
								status,
								type,
								type_id,
								amount,
								commission,
								currency_id,
								comment,
								comment_status
							)
							VALUES (
								'pending',
								'cf_project',
								?,
								?,
								?,
								?,
								?,
								'decrypted'
							)`, projectId, amount, commission, currencyId, commentText)
			if err != nil {
				return "", utils.ErrInfo(err)
			}
		}

		if len(comment) == 0 {
			comment = []byte("null");
		} else {
			comment = utils.HexToBin(comment);
		}

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(projectId)...)
		data = append(data, utils.EncodeLengthPlusData(amount)...)
		data = append(data, utils.EncodeLengthPlusData(commission)...)
		data = append(data, utils.EncodeLengthPlusData(comment)...)
		data = append(data, binSignatures...)

	case "cash_request_out" :

		toUserId := []byte(c.r.FormValue("to_user_id"));
		currencyId := []byte(c.r.FormValue("currency_id"));
		amount := []byte(c.r.FormValue("amount"));
		comment := utils.HexToBin([]byte(c.r.FormValue("comment")));
		commentText := []byte(c.r.FormValue("comment_text"));
		hashCode := []byte(c.r.FormValue("hash_code"));
		code := []byte(c.r.FormValue("code"));


		verifyData := map[string]string {c.r.FormValue("to_user_id"): "int", c.r.FormValue("currency_id"): "int", c.r.FormValue("amount"): "amount", c.r.FormValue("code"): "cash_code"}
		err := CheckInputData(verifyData)
		if err != nil {
			return "", utils.ErrInfo(err)
		}

		if c.SessRestricted == 0 {

			// пишем в личную таблу
			err = c.ExecSql(`INSERT INTO  `+c.MyPrefix+`my_cash_requests (
								to_user_id,
								currency_id,
								amount,
								comment,
								code
							)
							VALUES (
								?,
								?,
								?,
								?,
								?
							)`, toUserId, currencyId, amount, commentText, code)
			if err != nil {
				return "", utils.ErrInfo(err)
			}
			err = c.ExecSql(`INSERT INTO
							`+c.MyPrefix+`my_dc_transactions (
								status,
								type,
								type_id,
								to_user_id,
								amount,
								currency_id,
								comment,
								comment_status
							)
							VALUES (
								'pending',
								'cash_request',
								?,
								?,
								?,
								?,
								?,
								'decrypted'
							)`, userId, toUserId, amount, currencyId, commentText)
			if err != nil {
				return "", utils.ErrInfo(err)
			}
		}

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(toUserId)...)
		data = append(data, utils.EncodeLengthPlusData(amount)...)
		data = append(data, utils.EncodeLengthPlusData(comment)...)
		data = append(data, utils.EncodeLengthPlusData(currencyId)...)
		data = append(data, utils.EncodeLengthPlusData(hashCode)...)
		data = append(data, binSignatures...)


	case "cash_request_in" :

		cashRequestId := []byte(c.r.FormValue("cash_request_id"));
		code := []byte(c.r.FormValue("code"));

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(cashRequestId)...)
		data = append(data, utils.EncodeLengthPlusData(code)...)
		data = append(data, binSignatures...)



	case "abuses" :

		abuses := []byte(c.r.FormValue("abuses"));

		// проверим, не делал слал ли юзер абузы за последние сутки.
		// если слал - то выходим.
		num, err := c.Single(`
					SELECT time
					FROM log_time_abuses
					WHERE user_id = ?`, userId).Int64()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		if ( num > 0 ) {
			return "", utils.ErrInfo(err)
		}
		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(abuses)...)
		data = append(data, binSignatures...)


	case "admin_ban_miners" :

		usersIds := []byte(c.r.FormValue("users_ids"));

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(usersIds)...)
		data = append(data, binSignatures...)



	case "admin_unban_miners" :

		usersIds := []byte(c.r.FormValue("users_ids"));

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(usersIds)...)
		data = append(data, binSignatures...)



	case "admin_variables" :  // админ изменил variables

		variables := []byte(c.r.FormValue("variables"));

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(variables)...)
		data = append(data, binSignatures...)



	case "admin_spots" : // админ обновил набор точек для проверки лиц

		exampleSpots := []byte(c.r.FormValue("example_spots"));
		segments := []byte(c.r.FormValue("segments"));
		tolerances := []byte(c.r.FormValue("tolerances"));
		compatibility := []byte(c.r.FormValue("compatibility"));

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(exampleSpots)...)
		data = append(data, utils.EncodeLengthPlusData(segments)...)
		data = append(data, utils.EncodeLengthPlusData(tolerances)...)
		data = append(data, utils.EncodeLengthPlusData(compatibility)...)
		data = append(data, binSignatures...)



	case "admin_message" : // админ отправил alert message

		message := []byte(c.r.FormValue("message"));
		currencyList := []byte(c.r.FormValue("currency_list"));

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(message)...)
		data = append(data, utils.EncodeLengthPlusData(currencyList)...)
		data = append(data, binSignatures...)



	case "change_primary_key" :

		publicKey1 := []byte(c.r.FormValue("public_key_1"));
		publicKey2 := []byte(c.r.FormValue("public_key_2"));
		publicKey3 := []byte(c.r.FormValue("public_key_3"));
		privateKey := []byte(c.r.FormValue("private_key"));
		passwordHash := []byte(c.r.FormValue("password_hash"));
		savePrivateKey := utils.StrToInt(c.r.FormValue("save_private_key"))


		verifyData := map[string]string {c.r.FormValue("public_key_1"): "public_key", c.r.FormValue("password_hash"): "sha256"}
		err := CheckInputData(verifyData)
		if err != nil {
			return "", utils.ErrInfo(err)
		}

		if len(privateKey)>0 && !utils.CheckInputData(privateKey, "private_key") {
			return `incorrect private_key`, nil
		}

		if c.SessRestricted == 0 {
			if (savePrivateKey == 1 && c.Community == false) {
				err = c.ExecSql(`INSERT INTO  `+c.MyPrefix+`my_keys (
										public_key,
										private_key,
										password_hash
									)
									VALUES (
										[hex],
										?,
										?
									)`, publicKey1, privateKey, passwordHash)
				if err != nil {
					return "", utils.ErrInfo(err)
				}
			} else {
				err = c.ExecSql(`INSERT INTO  `+c.MyPrefix+`my_keys (
										public_key
									)
									VALUES (
										[hex]
									)`, publicKey1)
				if err != nil {
					return "", utils.ErrInfo(err)
				}
			}
		}

		bin_public_key_1 := utils.HexToBin(publicKey1);
		bin_public_key_2 := utils.HexToBin(publicKey2);
		bin_public_key_3 := utils.HexToBin(publicKey3);
		binPublicKeyPack :=  utils.EncodeLengthPlusData(bin_public_key_1)
		binPublicKeyPack = append(binPublicKeyPack, utils.EncodeLengthPlusData(bin_public_key_2)...)
		binPublicKeyPack = append(binPublicKeyPack, utils.EncodeLengthPlusData(bin_public_key_3)...)

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(binPublicKeyPack)...)
		data = append(data, binSignatures...)


	case "change_node_key" :

		publicKey := []byte(c.r.FormValue("public_key"));
		privateKey := []byte(c.r.FormValue("private_key"));

		verifyData := map[string]string {c.r.FormValue("public_key"): "public_key", c.r.FormValue("private_key"): "private_key"}
		err := CheckInputData(verifyData)
		if err != nil {
			return "", utils.ErrInfo(err)
		}

		if c.SessRestricted == 0 {
			err = c.ExecSql(`INSERT INTO  `+c.MyPrefix+`my_node_keys (
									public_key,
									private_key
								)
								VALUES (
									[hex],
									?
								)`, publicKey, privateKey)
			if err != nil {
				return "", utils.ErrInfo(err)
			}
		}

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(utils.HexToBin(publicKey))...)
		data = append(data, binSignatures...)



	case "votes_complex" :

		jsonData := []byte(c.r.FormValue("json_data"));

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(jsonData)...)
		data = append(data, binSignatures...)



	case "new_holidays" :

		startTime := []byte(c.r.FormValue("start_time"));
		endTime := []byte(c.r.FormValue("end_time"));

		verifyData := map[string]string {c.r.FormValue("start_time"): "int", c.r.FormValue("end_time"): "int"}
		err := CheckInputData(verifyData)
		if err != nil {
			return "", utils.ErrInfo(err)
		}

		if c.SessRestricted == 0 {
			err = c.ExecSql(`INSERT INTO
								`+c.MyPrefix+`my_holidays (
									start_time,
									end_time
								)
								VALUES (
									?,
									?
								)`, startTime, endTime)
			if err != nil {
				return "", utils.ErrInfo(err)
			}
		}

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(startTime)...)
		data = append(data, utils.EncodeLengthPlusData(endTime)...)
		data = append(data, binSignatures...)


	case "new_miner_update" :

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, binSignatures...)

		if c.SessRestricted == 0 {
			err = c.ExecSql(`UPDATE `+c.MyPrefix+`my_table
					SET node_voting_send_request = ?`, txTime)
			if err != nil {
				return "", utils.ErrInfo(err)
			}
		}



	case "admin_add_currency" :

		currencyName := []byte(c.r.FormValue("currency_name"));
		currencyFullName := []byte(c.r.FormValue("currency_full_name"));
		maxPromisedAmount := []byte(c.r.FormValue("max_promised_amount"));
		maxOtherCurrencies := []byte(c.r.FormValue("max_other_currencies"));

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(currencyName)...)
		data = append(data, utils.EncodeLengthPlusData(currencyFullName)...)
		data = append(data, utils.EncodeLengthPlusData(maxPromisedAmount)...)
		data = append(data, utils.EncodeLengthPlusData(maxOtherCurrencies)...)
		data = append(data, binSignatures...)



	case "admin_new_version" :

		softType := []byte(c.r.FormValue("soft_type"));
		version := []byte(c.r.FormValue("version"));
		format := []byte(c.r.FormValue("format"));

		newFile, err := ioutil.ReadFile("public/new.zip")
		if err != nil {
			return "", utils.ErrInfo(err)
		}

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(softType)...)
		data = append(data, utils.EncodeLengthPlusData(version)...)
		data = append(data, utils.EncodeLengthPlusData(newFile)...)
		data = append(data, utils.EncodeLengthPlusData(format)...)
		data = append(data, binSignatures...)



	case "admin_new_version_alert" :

		softType := []byte(c.r.FormValue("soft_type"));
		version := []byte(c.r.FormValue("version"));

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(softType)...)
		data = append(data, utils.EncodeLengthPlusData(version)...)
		data = append(data, binSignatures...)



	case "admin_blog" :

		title := []byte(c.r.FormValue("title"));
		message := []byte(c.r.FormValue("message"));

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(title)...)
		data = append(data, utils.EncodeLengthPlusData(message)...)
		data = append(data, binSignatures...)


	case "message_to_admin" :



		messageId := []byte(c.r.FormValue("message_id"));
		//parentId := []byte(c.r.FormValue("parent_id"));
		//subject := []byte(c.r.FormValue("subject"));
		//message := []byte(c.r.FormValue("message"));
		//messageType := []byte(c.r.FormValue("message_type"));
		//messageSubtype := []byte(c.r.FormValue("message_subtype"));
		encryptedMessage := []byte(c.r.FormValue("encrypted_message"));

		verifyData := map[string]string {c.r.FormValue("message_id"): "int", c.r.FormValue("encrypted_message"): "hex_message"}
		err := CheckInputData(verifyData)
		if err != nil {
			return "", utils.ErrInfo(err)
		}

		if c.SessRestricted == 0 {
			err = c.ExecSql(`UPDATE `+c.MyPrefix+`my_admin_messages
							SET  status = 'my_pending',
									encrypted = [hex]
							WHERE id = ?`, encryptedMessage, messageId)
			if err != nil {
				return "", utils.ErrInfo(err)
			}
		}

		encryptedMessage = utils.HexToBin(encryptedMessage);

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(encryptedMessage)...)
		data = append(data, binSignatures...)




	case "admin_answer" :

		//parentId := []byte(c.r.FormValue("parent_id"));
		//message := []byte(c.r.FormValue("message"));
		encryptedMessage := []byte(c.r.FormValue("encrypted_message"));
		toUserId := []byte(c.r.FormValue("to_user_id"));

		encryptedMessage = utils.HexToBin(encryptedMessage);

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(toUserId)...)
		data = append(data, utils.EncodeLengthPlusData(encryptedMessage)...)
		data = append(data, binSignatures...)


	case "change_host" :

		host := []byte(c.r.FormValue("host"));

		if !utils.CheckInputData(c.r.FormValue("host"), "host") {
			return `incorrect host`, nil
		}

		var community []int64
		if c.SessRestricted == 0 {

			node_admin_access, err := c.NodeAdminAccess(c.SessUserId, c.SessRestricted)
			if err != nil {
				return "", utils.ErrInfo(err)
			}
			if c.Community && node_admin_access {
				community, err = c.GetCommunityUsers();
				if err != nil {
					return "", utils.ErrInfo(err)
				}
			} else {
				community = []int64{c.SessUserId}
			}

			var myPrefix string
			for i:=0; i < len(community); i++ {
				if c.Community && node_admin_access {
					myPrefix = utils.Int64ToStr(community[i])+"_"
				} else if c.Community {
					myPrefix = c.MyPrefix
				} else {
					myPrefix = ""
				}
				uId := community[i]
				err = c.ExecSql(`
							UPDATE `+myPrefix+`my_table
							SET  host = ?,
									host_status = 'my_pending'`, host)
				if err != nil {
					return "", utils.ErrInfo(err)
				}
				nodePrivateKey, err := c.Single(`
							SELECT private_key
							FROM `+myPrefix+`my_node_keys
							WHERE block_id = (SELECT max(block_id) FROM `+myPrefix+`my_node_keys )`).Bytes()
				if err != nil {
					return "", utils.ErrInfo(err)
				}
				timeNow := time.Now().Unix()

				// подписываем нашим нод-ключем данные транзакции
				// Extract the PEM-encoded data block
				block, _ := pem.Decode(nodePrivateKey)
				if block == nil {
					return "", (errors.New("bad key data"))
				}
				if got, want := block.Type, "RSA PRIVATE KEY"; got != want {
					return "", (errors.New("unknown key type "+got+", want "+want));
				}
				privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
				if err != nil {
					return "", utils.ErrInfo(err)
				}
				forSign := fmt.Sprintf("%d,%d,%d,%s", utils.TypeInt(txType_), timeNow, uId, host)
				binSignature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, utils.HashSha1(forSign))
				if err != nil {
					return "", utils.ErrInfo(err)
				}

				// создаем новую транзакцию - подверждение, что фото скопировано и проверено.
				data = utils.DecToBin(txType, 1)
				data = append(data, utils.DecToBin(timeNow, 4)...)
				data = append(data, utils.EncodeLengthPlusData(utils.Int64ToByte(uId))...)
				data = append(data, utils.EncodeLengthPlusData(host)...)
				data = append(data, binSignature...)

				err = c.ExecSql("INSERT INTO queue_tx (hash, data) VALUES ([hex], [hex])", utils.Md5(data), utils.BinToHex(data))
				if err != nil {
					return "", utils.ErrInfo(err)
				}
			}
			return "ok", nil
		} else {
			return "access error", nil
		}

	case "new_forex_order" :

		sellCurrencyId := []byte(c.r.FormValue("sell_currency_id"));
		sellRate := []byte(c.r.FormValue("sell_rate"));
		amount := []byte(c.r.FormValue("amount"));
		buyCurrencyId := []byte(c.r.FormValue("buy_currency_id"));
		commission := []byte(c.r.FormValue("commission"));

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(sellCurrencyId)...)
		data = append(data, utils.EncodeLengthPlusData(sellRate)...)
		data = append(data, utils.EncodeLengthPlusData(amount)...)
		data = append(data, utils.EncodeLengthPlusData(buyCurrencyId)...)
		data = append(data, utils.EncodeLengthPlusData(commission)...)
		data = append(data, binSignatures...)



	case "for_repaid_fix" :

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, binSignatures...)



	case "actualization_promised_amounts" :

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, binSignatures...)

	case "change_commission" :

		commission := []byte(c.r.FormValue("commission"));
		commissionDecode := make(map[string][3]float64)
		err = json.Unmarshal(commission, &commissionDecode)
		if err != nil {
			return "", utils.ErrInfo(err)
		}

		poolCommission := make(map[string][3]float64)
		if c.Community {
			pool_commission_, err := c.Single(`
					SELECT commission
					FROM config`).Bytes()
			if err != nil {
				return "", utils.ErrInfo(err)
			}
			err = json.Unmarshal(pool_commission_, &poolCommission)
			if err != nil {
				return "", utils.ErrInfo(err)
			}
		}
		for currencyId, data := range commissionDecode {
			if !utils.CheckInputData(currencyId, "bigint") {
				return "", errors.New("incorrect currencyId")
			}
			// % от 0 до 10
			if !utils.CheckInputData(data[0], "currency_commission")  || data[0] > 10{
				return "", errors.New("incorrect pct")
			}
			// минимальная комиссия от 0. При 0% будет = 0
			if !utils.CheckInputData(data[1], "currency_commission") {
				return "", errors.New("incorrect currency_min_commission")
			}
			// макс. комиссия. 0 - значит, считается по %
			if !utils.CheckInputData(data[2], "currency_commission") {
				return "", errors.New("incorrect currency_max_commission")
			}
			if data[1] > data[2] && data[2]>0 {
				return "", errors.New("incorrect currency_max_commission")
			}
			// и если в пуле, то
			if len(poolCommission) > 0 {
				// нельзя допустить, чтобы блок подписал майнер, у которого комиссия больше той, что разрешана в пуле,
				// т.к. это приведет к попаднию в блок некорректной тр-ии, что приведет к сбою пула
				if len(poolCommission[currencyId]) > 0 && data[0] > poolCommission[currencyId][0] {
					return "", errors.New("incorrect commission")
				}
				if len(poolCommission[currencyId]) > 0 && data[1] > poolCommission[currencyId][1] {
					return "", errors.New("incorrect commission")
				}
			}
		}
		if c.SessRestricted == 0 {
			for currencyId, data := range commissionDecode {
				err = c.ExecSql(`
						INSERT INTO `+c.MyPrefix+`my_commission (
								currency_id,
								pct,
								min,
								max
							)
							VALUES (
								?,
								?,
								?,
								?
							)
	                    ON DUPLICATE KEY UPDATE pct=?, min=?, max=?`, currencyId, data[0], data[1], data[2], data[0], data[0], data[0])
				if err != nil {
					return "", utils.ErrInfo(err)
				}
			}
		}

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(commission)...)
		data = append(data, binSignatures...)


	case "change_key_active" :

		secret := utils.HexToBin([]byte(c.r.FormValue("secret")));

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(secret)...)
		data = append(data, binSignatures...)



	case "change_key_close" :

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, binSignatures...)



	case "change_key_request" :

		toUserId := []byte(c.r.FormValue("to_user_id"));

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(toUserId)...)
		data = append(data, binSignatures...)



	case "admin_change_primary_key" :

		forUserId := []byte(c.r.FormValue("for_user_id"));
		newPublicKey := utils.HexToBin([]byte(c.r.FormValue("new_public_key")));

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData(forUserId)...)
		data = append(data, utils.EncodeLengthPlusData(newPublicKey)...)
		data = append(data, binSignatures...)



	case "change_arbitrator_list" :

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData([]byte(c.r.FormValue("arbitration_trust_list")))...)
		data = append(data, binSignatures...)



	case "money_back_request" :

		var arbitratorEncText [5]string
		err := json.Unmarshal([]byte(c.r.PostFormValue("arbitrator_enc_text")), &arbitratorEncText)
		if err != nil {
			return fmt.Sprintf("%q", err), err
		}

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData([]byte(c.r.FormValue("order_id")))...)
		data = append(data, utils.EncodeLengthPlusData(utils.HexToBin([]byte(arbitratorEncText[0])))...)
		data = append(data, utils.EncodeLengthPlusData(utils.HexToBin([]byte(arbitratorEncText[1])))...)
		data = append(data, utils.EncodeLengthPlusData(utils.HexToBin([]byte(arbitratorEncText[2])))...)
		data = append(data, utils.EncodeLengthPlusData(utils.HexToBin([]byte(arbitratorEncText[3])))...)
		data = append(data, utils.EncodeLengthPlusData(utils.HexToBin([]byte(arbitratorEncText[4])))...)
		data = append(data, utils.EncodeLengthPlusData(utils.HexToBin([]byte(c.r.FormValue("seller_enc_text"))))...)
		data = append(data, binSignatures...)



	case "change_seller_hold_back" :

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData([]byte(c.r.FormValue("arbitration_days_refund")))...)
		data = append(data, utils.EncodeLengthPlusData([]byte(c.r.FormValue("hold_back_pct")))...)
		data = append(data, binSignatures...)



	case "money_back" :

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData([]byte(c.r.FormValue("order_id")))...)
		data = append(data, utils.EncodeLengthPlusData([]byte(c.r.FormValue("amount")))...)
		data = append(data, binSignatures...)



	case "change_arbitrator_conditions" :

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData([]byte(c.r.FormValue("conditions")))...)
		data = append(data, utils.EncodeLengthPlusData([]byte(c.r.FormValue("url")))...)
		data = append(data, binSignatures...)



	case "change_money_back_time" :

		data = utils.DecToBin(txType, 1)
		data = append(data, utils.DecToBin(txTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(userId)...)
		data = append(data, utils.EncodeLengthPlusData([]byte(c.r.FormValue("order_id")))...)
		data = append(data, utils.EncodeLengthPlusData([]byte(c.r.FormValue("days")))...)
		data = append(data, binSignatures...)


	}

	md5 := utils.Md5(data)
	if utils.InSliceString(txType_, []string{"new_pct", "new_max_promised_amounts", "new_reduction", "votes_node_new_miner", "new_max_other_currencies"}) {
		err := c.ExecSql(`INSERT INTO transactions_status (
				hash,
				time,
				type,
				user_id
			)
			VALUES (
				[hex],
				?,
				?,
				?
			)`, md5, time.Now().Unix(), txType, userId)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
	}
	err = c.ExecSql("INSERT INTO queue_tx (hash, data) VALUES ([hex], [hex])", md5, utils.BinToHex(data))
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	return `{"error":"null"}`, nil
}

func CheckInputData(data map[string]string) (error) {
	for k, v := range data {
		if !utils.CheckInputData(k, v) {
			return utils.ErrInfo(fmt.Errorf("incorrect "+v))
		}
	}
	return nil
}
