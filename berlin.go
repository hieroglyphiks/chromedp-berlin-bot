package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

const (
	desiredCountry          = "United States of America"
	countryNameAttributeIdx = 5
)

func bot(ctx context.Context) error {
	err := initializeSearch(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize search: %w", err)
	}

	<-ctx.Done()
	return nil
}

func initializeSearch(ctx context.Context) error {
	log.Println("Navigating to start page...")
	bookAppointmentButtonXPath := `//*[@id="mainForm"]/div/div/div/div/div/div/div/div/div/div[1]/div[1]/div[2]/a`

	chromedp.Run(ctx,
		network.ClearBrowserCookies(),
		network.ClearBrowserCache(),
		chromedp.Navigate("https://otv.verwalt-berlin.de/ams/TerminBuchen?lang=en"),
		chromedp.WaitVisible(bookAppointmentButtonXPath, chromedp.BySearch),
		chromedp.Click(bookAppointmentButtonXPath, chromedp.BySearch, chromedp.NodeVisible),
	)

	log.Println("Done.")

	log.Println("Accepting agreement...")
	agreementCheckboxXPath := `//*[@id="xi-cb-1"]`
	agreementButtonXPath := `//*[@id="applicationForm:managedForm:proceed"]`

	chromedp.Run(ctx,
		chromedp.WaitVisible(agreementCheckboxXPath, chromedp.BySearch),
		chromedp.Click(agreementCheckboxXPath, chromedp.BySearch, chromedp.NodeVisible),
		chromedp.Click(agreementButtonXPath, chromedp.BySearch, chromedp.NodeVisible),
	)

	log.Println("Done.")
	log.Println("Filling Out Form...")

	err := fillForm(ctx)
	if err != nil {
		return fmt.Errorf("failed to fill form: %w", err)
	}

	log.Println("Checking for unfilled form...")
	formBodyXPath := `//*[@id="xi-fs-19"]`
	formBodyHTML := ""

	chromedp.Run(ctx,
		chromedp.OuterHTML(formBodyXPath, &formBodyHTML, chromedp.BySearch),
	)

	for strings.Contains(formBodyHTML, "field is required") {
		log.Println("Form is not properly filled, refilling...")
		err := fillForm(ctx)
		if err != nil {
			return fmt.Errorf("failed to fill form: %w", err)
		}
	}

	log.Println("Done.")
	log.Println("Clicking Extend A Residence Title ...")
	extendAResidenceTitleXPath := `//*[@id="xi-div-30"]/div[2]/label/p`
	extendAResidenceTitleHTML := ""

	chromedp.Run(ctx,
		chromedp.WaitVisible(extendAResidenceTitleXPath, chromedp.BySearch),
		chromedp.OuterHTML(extendAResidenceTitleXPath, &extendAResidenceTitleHTML, chromedp.BySearch),
	)

	for !strings.Contains(extendAResidenceTitleHTML, "Extend a residence title") {
		fmt.Println("Selected Block Doesn't Contain Extend A Residence Title. Toggling Family Members...")
		fmt.Printf("Block HTML: %s\n", extendAResidenceTitleHTML)
		err := toggleFamilyMembers(ctx)
		if err != nil {
			return fmt.Errorf("failed to toggle family members: %w", err)
		}

		chromedp.Run(ctx,
			chromedp.WaitVisible(extendAResidenceTitleXPath, chromedp.BySearch),
			chromedp.OuterHTML(extendAResidenceTitleXPath, &extendAResidenceTitleHTML, chromedp.BySearch),
		)
	}

	chromedp.Run(ctx,
		chromedp.Click(extendAResidenceTitleXPath, chromedp.BySearch, chromedp.NodeVisible),
	)

	log.Println("Done.")
	log.Println("Clicking Economic Activity...")
	economicActivityXPath := `//*[@id="inner-368-0-2"]/div/div[3]/label`

	chromedp.Run(ctx,
		chromedp.WaitVisible(economicActivityXPath, chromedp.BySearch),
		chromedp.Click(economicActivityXPath, chromedp.BySearch, chromedp.NodeVisible),
	)

	log.Println("Done.")
	log.Println("Clicking Academic Education...")
	academicEducationXPath := `//*[@id="inner-368-0-2"]/div/div[4]/div/div[5]/label`
	infoAboutSelectedServiceXPath := `//*[@id="xi-fs-28"]/legend`

	chromedp.Run(ctx,
		chromedp.WaitVisible(academicEducationXPath, chromedp.BySearch),
		chromedp.Click(academicEducationXPath, chromedp.BySearch, chromedp.NodeVisible),
		chromedp.WaitVisible(infoAboutSelectedServiceXPath, chromedp.BySearch),
	)

	log.Println("Done.")
	log.Println("Searching For an Appointment...")
	searchAppointmentXPath := `//*[@id="applicationForm:managedForm:proceed"]`

	chromedp.Run(ctx,
		chromedp.WaitVisible(searchAppointmentXPath, chromedp.BySearch),
		chromedp.Click(searchAppointmentXPath, chromedp.BySearch, chromedp.NodeVisible),
		chromedp.WaitVisible(searchAppointmentXPath, chromedp.BySearch),
	)

	log.Println("Done.")
	t := time.Now()
	for time.Since(t) < 20*time.Minute {
		log.Println("Checking for error message...")

		var html string
		chromedp.Run(ctx,
			// chromedp.Nodes(`//*[@id="messagesBox"]/ul/li`, &errors, chromedp.BySearch),
			chromedp.OuterHTML("#main", &html, chromedp.ByQuery),
		)

		if strings.Contains(html, "Time Selection") {
			fmt.Println(html)
			fmt.Println("WTF FOUND AN APPOINTMENT")
			break
		} else {
			fmt.Println("No dates found, retrying")
		}

		fmt.Println("Done.")

		chromedp.Run(ctx,
			chromedp.Sleep(10*time.Second),
			chromedp.WaitVisible(searchAppointmentXPath, chromedp.BySearch),
			chromedp.Click(searchAppointmentXPath, chromedp.BySearch, chromedp.NodeVisible),
			chromedp.WaitVisible(searchAppointmentXPath, chromedp.BySearch),
		)
		fmt.Println("Done.")
	}

	return nil
}

func fillForm(ctx context.Context) error {
	countrySelectorDropdownXPath := `//*[@id="xi-sel-400"]`
	countries := []*cdp.Node{}

	chromedp.Run(ctx,
		chromedp.WaitVisible(countrySelectorDropdownXPath, chromedp.BySearch),
		chromedp.Nodes(countrySelectorDropdownXPath, &countries, chromedp.BySearch),
	)

	for _, c := range countries {
		for _, child := range c.Children {
			if len(child.Attributes) < countryNameAttributeIdx {
				fmt.Printf("%+v\n", child)
				return fmt.Errorf("malformed country dropdown attributes")
			}

			if strings.Contains(child.Attributes[5], desiredCountry) {
				err := chromedp.Run(ctx,
					// this form is super janky idk why but it needs this
					chromedp.Sleep(2*time.Second),
					chromedp.SetValue(`//select[@id="xi-sel-400"]`, child.Attributes[3], chromedp.BySearch),

					// this sets it to one person
					chromedp.Sleep(time.Second),
					chromedp.SetValue(`//select[@id="xi-sel-422"]`, "1", chromedp.BySearch),
					// this sets it to no other family
					chromedp.Sleep(time.Second),
					chromedp.SetValue(`//select[@id="xi-sel-427"]`, "2", chromedp.BySearch),
				)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}

	return nil
}

func toggleFamilyMembers(ctx context.Context) error {
	err := chromedp.Run(ctx,
		// this sets it to yes family
		chromedp.Sleep(time.Second),
		chromedp.SetValue(`//select[@id="xi-sel-427"]`, "1", chromedp.BySearch),

		// this sets it to no other family
		chromedp.Sleep(time.Second),
		chromedp.SetValue(`//select[@id="xi-sel-427"]`, "2", chromedp.BySearch),
	)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}
