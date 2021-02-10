package main

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/jlaffaye/ftp"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	a := app.New()
	w := a.NewWindow("HTTP/FTP client")
	hello := widget.NewLabel("Select client")
	ftp := widget.NewButton("FTP", func() {
		runFtp(w)
	})
	http := widget.NewButton("HTTP", func() {
		runHttp(w)
	})
	w.SetContent(widget.NewVBox(
		hello,
		ftp,
		http,
	))
	w.Resize(fyne.NewSize(400,100))
	w.ShowAndRun()
}
func runFtp(w fyne.Window){
	var emptyValidator fyne.StringValidator = func(str string) error {
		if str == ""{
			return fmt.Errorf("Error: empty string!")
		}
		return nil
	}
	var validator fyne.StringValidator = func(str string) error{
		if !strings.HasPrefix(str,"ftp"){
			return fmt.Errorf("wrong format")
		}
		return nil
	}
	var selector = new(widget.SelectEntry)
	input := widget.NewEntry()
	pass := widget.NewPasswordEntry()
	user := widget.NewEntry()

	input.SetText("ftp.dlptest.com:21")
	user.SetText("dlpuser@dlptest.com")
	pass.SetText("eUj8GeW55SvYaswqUyDSm5v6N")
	text := widget.NewLabel("")

	button := widget.NewButton("Get ftp content", func() {
		input.Validator = validator
		pass.Validator = emptyValidator
		user.Validator = emptyValidator
		err := func() error {
			err := input.Validate()
			if err != nil{
				text.SetText(err.Error())
				return err
			}
			err = pass.Validate()
			if err != nil{
				text.SetText(err.Error()+"in pass")
				return err
			}
			err = user.Validate()
			if err != nil{
				text.SetText(err.Error()+"in user")
				return err
			}
			return nil
		}()
		if err != nil{
			input.SetText("")
		}else{
			var item string
			c, err := ftp.Dial(input.Text, ftp.DialWithTimeout(5*time.Second))
			if err != nil {
				log.Fatal(err)
			}

			err = c.Login(user.Text, pass.Text)
			if err != nil {
				log.Fatal(err)
			}
			list, err := c.List(".")
			text.SetText("")
			listing := make([]string, 1)

			for _, v := range list{
				text.Text = text.Text + v.Name + "\n" + strconv.Itoa(int(v.Size)) + " bytes\n"
				listing = append(listing, v.Name)
			}
			selector = widget.NewSelectEntry(listing)
			scroller := widget.NewVScrollContainer(text)
			pbar := widget.NewProgressBarInfinite()
			pbar.Hide()
			fpath := widget.NewEntry()
			selector.OnChanged = func(s string) {
				fmt.Println(s)
			 	item = s
			}

			scroller.SetMinSize(fyne.NewSize(300,500))
			w.SetContent(widget.NewVBox(
				selector,
				scroller,
				fpath,
				pbar,
				widget.NewButton("Download", func() {
					pbar.Show()
					var fp string
					res, _ := c.Retr(item)
					if err != nil {
						log.Fatal(err)
					}
					defer res.Close()
					if fpath.Text != ""{
						fp = fpath.Text
					}else{
						fp = item
					}
					outFile, err := os.Create(fp)
					if err != nil {
						log.Fatal(err)
					}
					defer outFile.Close()
					_, err = io.Copy(outFile, res)
					if err != nil {
						log.Fatal(err)
					}
					pbar.Hide()
				}),
				widget.NewButton("Disconnect", func() {
					c.Quit()
				}),
			))
		}
	})
	w.SetContent(widget.NewVBox(
		user,
		pass,
		input,
		button,
	))
}
func runHttp(w fyne.Window){
	var validator fyne.StringValidator = func(str string) error{
		if !strings.HasPrefix(str,"http://") && !strings.HasPrefix(str,"https://"){
			return fmt.Errorf("wrong format")
		}
		return nil
	}
	input := widget.NewEntry()
	text := widget.NewLabel("")
	button := widget.NewButton("Get http content", func() {
		input.Validator = validator
		err := input.Validate()
		if err != nil{
			text.SetText("Wrong format! input url with http:// or https:// prefix")
			input.SetText("")
		}else{
			resp, err := http.Get(input.Text)
			if err != nil{
				text.SetText("An error occurred:"+ err.Error())
				input.SetText("")
			}else{
				responseData,err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Fatal(err)
				}
				text.SetText(string(responseData))
				text.Resize(fyne.NewSize(300,300))
			}
		}
	})
	scroller := widget.NewVScrollContainer(text)
	scroller.SetMinSize(fyne.NewSize(300,300))
	w.Resize(fyne.NewSize(400, 400))
	w.SetContent(widget.NewVBox(
		scroller,
		input,
		button,
		))
}