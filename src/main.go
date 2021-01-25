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
	var numbEx = 0
	var exercise = 1
	var command string
	var btnPressed bool
	var goOutOrStay = false
	var firstEx = false

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
			case "дальше":
				r.Request.Command = "дальше"
			case "стоп":
				r.Request.Command = marusia.OnInterrupt
			case "продолжить":
				if firstEx == true {
					r.Request.Command = "дальше"
					firstEx = false
				} else {
					r.Request.Command = command
				}
				goOutOrStay = false
			}
			btnPressed = true
			returnToTheBeginning = true
		}

		if r.Request.Type == marusia.SimpleUtterance || btnPressed {
			for returnToTheBeginning {
				switch r.Request.Command {
				case marusia.OnStart:
					resp.Text = "Гимнастика для глаз позволит Вам сохранить зрение, а иногда и улучшить его. Главное – заниматься каждый день. У меня подготовлены для Вас комплексы упражнений длиной 5, 10 и 15 минут. Какой из них Вы хотите выбрать?"
					resp.TTS = "Гимнастика для глаз позволит Вам сохранить зрение, а иногда и улучшить его. Главное – заниматься каждый день. У меня подготовлены для Вас комплексы упражнений длиной пять, десять и пятнадцать минут. Какой из них Вы хотите выбрать?"
					resp.Card = marusia.NewBigImage(
						"",
						"",
						457239088,
					)
					resp = btnComplexOfExercise(resp)
					numbEx = 0
					exercise = 1
					command = ""
					returnToTheBeginning = false
				case "да", "конечно", "хочу":
					if !goOutOrStay {
						r.Request.Command = "default"
					} else {
						resp.Text = "Приятно было с Вами пообщаться. Возвращайтесь, когда Вам будет удобно."
						goOutOrStay = false
						returnToTheBeginning = false
						resp.EndSession = true
					}
				case "нет", "не хочу":
					if !goOutOrStay {
						r.Request.Command = "default"
					} else {
						if firstEx == true {
							r.Request.Command = "дальше"
							firstEx = false
						} else {
							r.Request.Command = command
						}
						goOutOrStay = false
					}
				case "5", "пять", "5 минут", "пять минут":
					resp = firstExercise(resp)
					firstEx = true
					resp = btnExecuteAndStop(resp)
					numbEx = 8
					exercise = 1
					command = "5"
					returnToTheBeginning = false
				case "10", "десять", "10 минут", "десять минут":
					resp = firstExercise(resp)
					firstEx = true
					resp = btnExecuteAndStop(resp)
					numbEx = 13
					exercise = 1
					command = "10"
					returnToTheBeginning = false
				case "15", "пятнадцать", "15 минут", "пятнадцать минут":
					resp = firstExercise(resp)
					firstEx = true
					resp = btnExecuteAndStop(resp)
					numbEx = 19
					exercise = 1
					command = "15"
					returnToTheBeginning = false
				case "продолжаем", "продолжим", "продолжить",
						"далее", "дальше", "готово",
						"выполнено", "выполнена", "выполнил", "выполнила",
						"закончил", "закончила", "закончено",
						"все", "сделано", "сделал", "сделала" :
					if exercise <= numbEx && numbEx != 0 {
						resp = nextExercise(exercise, resp)
						exercise++
						resp = btnExecuteAndStop(resp)
						command = "дальше"
						returnToTheBeginning = false
					} else if numbEx == 0 && !goOutOrStay {
						r.Request.Command = "default"
					} else if exercise > numbEx {
						resp.Text = "На этом все. И не забывайте, что для достижения наилучшего эффекта, необходимо выполнять гимнастику каждый день."
						returnToTheBeginning = false
						resp.EndSession = true
					}
				case marusia.OnInterrupt:
					resp.Text = "Приятно было с Вами пообщаться. Возвращайтесь, когда Вам будет удобно."
					goOutOrStay = false
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
					goOutOrStay = true
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
	resp.AddButton("Дальше", MyPayload {
		Text: "дальше",
	})
	resp.AddButton("Стоп", MyPayload {
		Text: "стоп",
	})
	return resp
}

func firstExercise(resp marusia.Response) marusia.Response {
	resp.Text = "Вот Ваше первое упражнение. После выполнения очередного упражнения скажите \"Дальше\" - для продолжения или \"Стоп\" - для завершения. \nКрепко зажмурьте глаза на 30 секунд."
	resp.TTS = "Вот Ваше первое упражнение. После выполнения очередного упражнения скажите \"Дальше\" - для продолжения или \"Стоп\" - для завершения. Крепко зажмурьте глаза на тридцать секунд. <speaker audio_vk_id=\"-2000512015_456239053\">. ... Время вышло."
	resp.Card = marusia.NewBigImage(
		"",
		"",
		457239041,
	)
	return resp
}

func nextExercise(exercise int, resp marusia.Response) marusia.Response {
	switch exercise {
	case 1:
		resp.Text = "Медленно посмотрите слева направо и справа налево. Выполняйте несколько раз."
		resp.TTS = "Медленно посмотрите слева направо и справа налево. Выполняйте несколько раз. <speaker audio_vk_id=\"-2000512015_456239054\">"
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239042,
		)
	case 2:
		resp.Text = "Медленно посмотрите слева направо по диагонали и справа налево по диагонали. Повторите несколько раз."
		resp.TTS = "Медленно посмотрите слева направо по диагонали и справа налево по диагонали. Повторите несколько раз. <speaker audio_vk_id=\"-2000512015_456239055\">"
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239043,
		)
	case 3:
		resp.Text = "Медленно рисуйте глазами цифру восемь несколько раз."
		resp.TTS = "Медленно рисуйте глазами цифру восемь несколько раз. <speaker audio_vk_id=\"-2000512015_456239056\">"
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239044,
		)
	case 4:
		resp.Text = "Медленно рисуйте глазами большой круг несколько раз."
		resp.TTS = "Медленно рисуйте глазами большой круг несколько раз. <speaker audio_vk_id=\"-2000512015_456239057\">"
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239045,
		)
	case 5:
		resp.Text = "Смотрите между бровей на протяжении 20 секунд."
		resp.TTS = "Смотрите между бровей на протяжении двадцати секунд. <speaker audio_vk_id=\"-2000512015_456239058\">. ... Время прошло."
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239046,
		)
	case 6:
		resp.Text = "Смотрите на кончик носа на протяжении 20 секунд."
		resp.TTS = "Смотрите на кончик носа на протяжении двадцати секунд. <speaker audio_vk_id=\"-2000512015_456239059\">. ... Время вышло."
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239047,
		)
	case 7:
		resp.Text = "Смотрите вдаль около 20 секунд."
		resp.TTS = "Смотрите вдаль около двадцати секунд. <speaker audio_vk_id=\"-2000512015_456239060\">. ... Все! Двигаемся дальше."
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239086,
		)
	case 8:
		resp.Text = "Быстро и легко моргайте примерно 30 секунд."
		resp.TTS = "Быстро и легко моргайте примерно тридцать секунд. <speaker audio_vk_id=\"-2000512015_456239061\">. ... Все! Время прошло."
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239049,
		)
	case 9:
		resp.Text = "Разотрите ладони до тепла и прикройте ими глаза, скрестив пальцы на середине лба, так чтобы не сдавливались глаза и не просачивался свет. Постарайтесь расслабиться и представить что-нибудь приятное. Выполняйте 2 минуты."
		resp.TTS = "Разотрите ладони до тепла и прикройте ими глаза, скрестив пальцы на середине лба, так чтобы не сдавливались глаза и не просачивался свет. Постарайтесь расслабиться и представить что-нибудь приятное. Выполняйте две минуты.  <speaker audio_vk_id=\"-2000512015_456239062\">. ... Время вышло."
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239092,
		)
	case 10:
		resp.Text = "Слегка массируйте двумя пальцами каждой руки в области бровей от переносицы до висков около 30 секунд."
		resp.TTS = "Слегка массируйте двумя пальцами каждой руки в области бровей от переносицы до висков около тридцати секунд. <speaker audio_vk_id=\"-2000512015_456239063\">. ... Время истекло."
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239050,
		)
	case 11:
		resp.Text = "Слегка массируйте двумя пальцами каждой руки в области под глазами примерно 30 секунд."
		resp.TTS = "Слегка массируйте двумя пальцами каждой руки в области под глазами примерно тридцать секунд. <speaker audio_vk_id=\"-2000512015_456239064\">. ... Стоп! Двигаемся дальше."
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239051,
		)
	case 12:
		resp.Text = "Слегка нажмите тремя пальцами каждой руки на верхние веки, через 2 секунды снимите пальцы с век. Повторите 5 раз."
		resp.TTS = "Слегка нажмите тремя пальцами каждой руки на верхние веки, через две секунды снимите пальцы с век. Повторите пять раз. <speaker audio_vk_id=\"-2000512015_456239065\">"
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239052,
		)
	case 13:
		resp.Text = "Поднимите брови, а после опустите и нахмурьте их. Выполните 10 раз."
		resp.TTS = "Поднимите брови, а после опустите и нахмурьте их. Выполните десять раз. <speaker audio_vk_id=\"-2000512015_456239066\">"
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239053,
		)
	case 14:
		resp.Text = "Поднесите палец к переносице, сфокусируйтесь на нем и медленно отдаляйте палец от глаз, при этом продолжая на нем фокусироваться. Выполните 3 повторения."
		resp.TTS = "Поднесите палец к переносице, сфокусируйтесь на нем и медленно отдаляйте палец от глаз, при этом продолжая на нем фокусироваться. Выполните три повторения. <speaker audio_vk_id=\"-2000512015_456239067\">"
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239089,
		)
	case 15:
		resp.Text = "Поднесите палец к переносице, сфокусируйтесь на нем на 3 секунды и резко переведите взгляд на любой объект вдалеке также на 3 секунды. Выполните 5 повторений."
		resp.TTS = "Поднесите палец к переносице, сфокусируйтесь на нем на три секунды и резко переведите взгляд на любой объект вдалеке также на три секунды. Выполните пять повторений. <speaker audio_vk_id=\"-2000512015_456239068\">"
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239087,
		)
	case 16:
		resp.Text = "Подойдите к окну и начните рассматривать объекты вблизи и вдали в течении 30 секунд."
		resp.TTS = "Подойдите к окну и начните рассматривать объекты вблизи и вдали в течении тридцати секунд. <speaker audio_vk_id=\"-2000512015_456239069\">. ... Время вышло."
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239054,
		)
	case 17:
		resp.Text = "Прикройте рукой левый глаз и продолжайте рассматривать объекты на протяжении 30 секунд."
		resp.TTS = "Прикройте рукой левый глаз и продолжайте рассматривать объекты на протяжении тридцати секунд. <speaker audio_vk_id=\"-2000512015_456239070\">. ... Тридцать секунд прошло."
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239090,
		)
	case 18:
		resp.Text = "Прикройте рукой правый глаз и продолжайте рассматривать объекты около 30 секунд."
		resp.TTS = "Прикройте рукой правый глаз и продолжайте рассматривать объекты около тридцати секунд. <speaker audio_vk_id=\"-2000512015_456239070\">. ... Все! Идем дальше."
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239091,
		)
	case 19:
		resp.Text = "Слегка проморгайтесь и отдохните."
		resp.Card = marusia.NewBigImage(
			"",
			"",
			457239057,
		)
	}

	return resp
}

