package controllers

import (
    "helloapp/app/models"
    "github.com/revel/revel"
    "encoding/json"
)

type BidItemCtrl struct {
    GorpController
}

func (c BidItemCtrl) parseBidItem() (models.BidItem, error) {
    biditem := models.BidItem{}
    err := json.NewDecoder(c.Request.Body).Decode(&biditem)
    return biditem, err
}

func (c BidItemCtrl) Add() revel.Result {
    if biditem, err := c.parseBidItem(); err != nil {
        return c.RenderText("Unable to parse the BidItem from JSON.")
    } else {
        // Validate the model
        biditem.Validate(c.Validation)
        if c.Validation.HasErrors() {
            // Do something better here!
            return c.RenderText("You have error in your BidItem.")
        } else {
            if err := c.Txn.Insert(&biditem); err != nil {
                return c.RenderText(
                    "Error inserting record into database!")
            } else {
                return c.RenderJson(biditem)
            }
        }
    }
}

func (c BidItemCtrl) Get(id int64) revel.Result {
    biditem := new(models.BidItem)
    err := c.Txn.SelectOne(biditem, 
        `SELECT * FROM BidItem WHERE id = ?`, id)
    if err != nil {
        return c.RenderText("Error.  Item probably doesn't exist.")
    }
    return c.RenderJson(biditem)
}

func (c BidItemCtrl) List() revel.Result {
    lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
    limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
    biditems, err := c.Txn.Select(models.BidItem{}, 
        `SELECT * FROM BidItem WHERE Id > ? LIMIT ?`, lastId, limit)
    if err != nil {
        return c.RenderText(
            "Error trying to get records from DB.")
    }
    return c.RenderJson(biditems)
}

func (c BidItemCtrl) Update(id int64) revel.Result {
    biditem, err := c.parseBidItem()
    if err != nil {
        return c.RenderText("Unable to parse the BidItem from JSON.")
    }
    // Ensure the Id is set.
    biditem.Id = id
    success, err := c.Txn.Update(&biditem)
    if err != nil || success == 0 {
        return c.RenderText("Unable to update bid item.")
    }
    return c.RenderText("Updated %v", id)
}

func (c BidItemCtrl) Delete(id int64) revel.Result {
    success, err := c.Txn.Delete(&models.BidItem{Id: id})
    if err != nil || success == 0 {
        return c.RenderText("Failed to remove BidItem")
    }
    return c.RenderText("Deleted %v", id)
}
