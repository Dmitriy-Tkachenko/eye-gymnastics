package main

import (
"encoding/json"
"github.com/SevereCloud/vksdk/v2/marusia"
"net/http"
"os"
)

type MyPayload struct {
	Text string
}


func main() {
	wh := marusia.NewWebhook()
	var firstEntry = true
	var numbEx = 0
	var exercise = 1
	var command string
	var btnPressed bool
	var def = false

	wh.OnEvent(func(r marusia.Request) (resp marusia.Response) {
		var returnToTheBeginning = true
		if r.Request.Type == marusia.ButtonPressed {
			var p MyPayload
			err := json.Unmarshal(r.Request.Payload, &p)
			if err != nil {
				resp.Text = "Что-то пошло не так"
				return
			}

			switch p.Text {
			case "":
				r.Request.Command = marusia.OnStart
			case "да":
				r.Request.Command = "да"
			case "нет":
				r.Request.Command = "нет"
			case "5":
				r.Request.Command = "5"
			case "10":
				r.Request.Command = "10"
			case "15":
				r.Request.Command = "15"
			case "выполнено":
				r.Request.Command = "выполнено"
			case "стоп":
				r.Request.Command = marusia.OnInterrupt
			case "продолжить":
				if exercise != 1 && command == "выполнено" {
					exercise--
				}
				r.Request.Command = command
			}
			btnPressed = true
			returnToTheBeginning = true
		}

		if r.Request.Type == marusia.SimpleUtterance || btnPressed {
			for returnToTheBeginning {
				switch r.Request.Command {
				case marusia.OnStart:
					if firstEntry {
						resp.Text = "Гимнастика для глаз позволит вам сохранить зрение, а иногда и улучшить его. Главное - заниматься каждый день. Хотите начать?"
						resp.AddButton("Да", MyPayload{
							Text: "да",
						})
						resp.AddButton("Нет", MyPayload{
							Text: "нет",
						})
					} else {
						resp.Text = "Рада вас снова видеть! У меня подготовлены для вас комплексы упражнений длиной 5, 10 и 15 минут. Какой из них вы хотите выбрать?"
						resp = btnComplexOfExercise(resp)
						numbEx = 0
						exercise = 1
					}
					command = ""

					returnToTheBeginning = false
				case "да", "конечно", "хочу":
					if firstEntry || btnPressed {
						resp.Text = "У меня подготовлены для вас комплексы упражнений длиной 5, 10 и 15 минут. Какой из них вы хотите выбрать?"
						resp = btnComplexOfExercise(resp)
						firstEntry = false
						command = "да"
					} else if command == "default" {
						resp.Text = "Приятно было с вами пообщаться. Возвращайтесь, когда будет удобно."
						resp.EndSession = true
					} else {
						resp.Text = "Хотите завершить упражнения?"
						resp.AddButton("Да", MyPayload {
							Text: "стоп",
						})
						resp.AddButton("Нет", MyPayload {
							Text: "продолжить",
						})
					}
					returnToTheBeginning = false
				case "нет", "не хочу":
					if firstEntry || btnPressed {
						resp.Text = "Хорошо. Возвращайтесь, когда сможете."
						resp.EndSession = true
						command = "нет"
						returnToTheBeginning = false
					} else if command == "default" {
						r.Request.Command = command
						def = false
					} else {
						resp.Text = "Хотите завершить упражнения?"
						resp.AddButton("Да", MyPayload {
							Text: "стоп",
						})
						resp.AddButton("Нет", MyPayload {
							Text: "продолжить",
						})
						returnToTheBeginning = false
					}
				case "5", "пять", "5 минут", "пять минут":
					if !firstEntry {
						resp = firstExercise(resp)
						resp = btnExecuteAndStop(resp)
						numbEx = 6
						command = "5"
					}
					returnToTheBeginning = false
				case "10", "десять", "10 минут", "десять минут":
					if !firstEntry {
						resp = firstExercise(resp)
						resp = btnExecuteAndStop(resp)
						numbEx = 13
						command = "10"
					}
					returnToTheBeginning = false
				case "15", "пятнадцать", "15 минут", "пятнадцать минут":
					if !firstEntry {
						resp = firstExercise(resp)
						resp = btnExecuteAndStop(resp)
						numbEx = 19
						command = "15"
					}
					returnToTheBeginning = false
				case "продолжаем", "продолжим", "далее", "дальше", "готово", "выполнено", "все", "сделано":
					if exercise <= numbEx && numbEx != 0 {
						resp = exercises(exercise, resp)
						exercise++
						resp = btnExecuteAndStop(resp)
						command = "выполнено"
					} else if numbEx == 0 && !firstEntry {
						resp.Text = "Пожалуйста, выберите комплекс упражнений: 5, 10 или 15 минут?"
					} else if numbEx == 0 && firstEntry {
						resp.Text = "Вы готовы начать?"
					} else if exercise > numbEx {
						resp.Text = "На этом все. Возвращайтесь, когда сможете."
						resp.EndSession = true
					}
					returnToTheBeginning = false
				case marusia.OnInterrupt:
					resp.Text = "Приятно было с вами пообщаться. Возвращайтесь, когда будет удобно."
					resp.EndSession = true
					returnToTheBeginning = false
				default:
					resp.Text = "Хотите завершить упражнения?"
					resp.AddButton("Да", MyPayload {
						Text: "стоп",
					})
					resp.AddButton("Нет", MyPayload {
						Text: "продолжить",
					})
					def = true
					returnToTheBeginning = false
				}
				btnPressed = false
			}
		}
		return
	})
	http.HandleFunc("/",wh.HandleFunc)
	http.ListenAndServe(":" + os.Getenv("PORT"),nil)
}

func btnComplexOfExercise(resp marusia.Response) marusia.Response {
	resp.AddButton("5 мин.", MyPayload{
		Text: "5",
	})
	resp.AddButton("10 мин.", MyPayload{
		Text: "10",
	})
	resp.AddButton("15 мин.", MyPayload{
		Text: "15",
	})
	return resp
}

func btnExecuteAndStop(resp marusia.Response) marusia.Response {
	resp.AddButton("Выполнено", MyPayload {
		Text: "выполнено",
	})
	resp.AddButton("Стоп", MyPayload {
		Text: "стоп",
	})
	return resp
}

func firstExercise(resp marusia.Response) marusia.Response {
	resp.Text = "Вот ваше первое упражнение, после выполнения очередного упражнения просто скажите \"Выполнено\" - для продолжения или \"Стоп\" - для завершения. \nКрепко зажмурьте глаза на 30 секунд."
	resp.Card = marusia.NewBigImage(
		"",
		"",
		457239041,
	)
	return resp
}

func exercises(exercise int, resp marusia.Response) marusia.Response {
	switch exercise {
	case 1:
		resp.Text = "Медленно посмотрите слева на право и справа налево. Выполняйте несколько раз."
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239042,
		)
	case 2:
		resp.Text = "Медленно посмотрите слева направо по диагонали и справа на лево по диагонали. Повторите несколько раз."
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239043,
		)
	case 3:
		resp.Text = "Медленно рисуйте глазами цифру восемь несколько раз."
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239044,
		)
	case 4:
		resp.Text = "Медленно рисуйте глазами большой круг."
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239045,
		)
	case 5:
		resp.Text = "Смотрите между бровей на протяжении 20 секунд."
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239046,
		)
	case 6:
		resp.Text = "Смотрите на кончик носа на протяжении 20 секунд."
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239047,
		)
	case 7:
		resp.Text = "Смотрите вдаль около 20 секунд."
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239048,
		)
	case 8:
		resp.Text = "Быстро и легко моргайте примерно 30 секунд."
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239049,
		)
	case 9:
		resp.Text = "Разотрите ладони до тепла и прикройте ими глаза, скрестив пальцы на середине лба, так чтобы не сдавливались глаза и не просачивался свет. Постарайтесь расслабиться и представить что-нибудь приятное. Выполняйте 3 минуты."
	case 10:
		resp.Text = "Слегка массируйте двумя пальцами каждой руки в области бровей от переносицы до висков около 30 секунд."
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239050,
		)
	case 11:
		resp.Text = "Слегка массируйте двумя пальцами каждой руки в области под глазами примерно 30 секунд"
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239051,
		)
	case 12:
		resp.Text = "Слегка нажмите тремя пальцами каждой руки на верхние веки, через 2 секунды снимите пальцы с век. Повторите 5 раз."
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239052,
		)
	case 13:
		resp.Text = "Поднимите брови, а после опустите и нахмурьте их. Выполните 10 раз."
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239053,
		)
	case 14:
		resp.Text = "Поднесите палец к переносице, сфокусируйтесь на нем и медленно отдаляйте палец от глаз, при этом продолжая на нем фокусироваться. Выполните 3 повторения."
	case 15:
		resp.Text = "Поднесите палец к переносице, сфокусируйтесь на нем на 3 секунды и резко переведите взгляд на любой объект вдалеке также на 3 секунды. Выполните 5 повторений."
	case 16:
		resp.Text = "Подойдите к окну и начните рассматривать объекты вблизи и вдали в течении 30 секунд."
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239054,
		)
	case 17:
		resp.Text = "Прикройте рукой левый глаз и продолжайте рассматривать объекты на протяжении 30 секунд."
	case 18:
		resp.Text = "Прикройте рукой правый глаз и продолжайте рассматривать объекты около 30 секунд."
	case 19:
		resp.Text = "Слегка проморгайтесь и отдохните"
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239055,
		)
	}

	return resp
}

