package scraper

import (
	"encoding/json"
	"testing"
)

func BenchmarkHtmlTest(b *testing.B) {
	data := `<div xmlns:xs="http://www.w3.org/2001/XMLSchema" class="single-package"><input type="hidden" value="13559300030524" class="js-waybill"><input type="hidden" value="13559300030524" class="js-waybill-paczki"><fieldset class="compact">
      <div class="form-group"><span class="label">Przesyłka</span><span class="input"><span class="input-text">13559300030524</span></span></div>
      <div class="form-group"><label>Paczki w przesyłce</label><span class="input"><select name="parcel" class="custom-select">
               <option value="13559300030524">13559300030524</option></select></span></div>
      <div class="form-group"><span class="label">Paczka</span><span class="input"><span class="input-text">13559300030524&nbsp;
               </span></span></div>
      <div class="form-group"><span class="label">Jesteś odbiorcą? Możesz samodzielnie zarządzać przesyłką</span><span class="input"><a href="https://mojapaczka.dpd.com.pl/login?parcel=13559300030524">
               <div class="column arrow"><span class="btn--arrow-right btn--no-text"></span></div></a></span><span class="label l-info"><span class="l-info">Portal pozwala na przekierowanie, zmianę daty doręczenia, rezygnację z paczki i przekierowanie do punktu Pickup</span></span></div>
      <div class="js-package-details subform" style="display: none;"></div>
   </fieldset>
   <fieldset class="compact">
      <h3>Historia przesyłki</h3>
      <div class="table-wrapper-400">
         <table class="table-track">
            <thead>
               <th>Data</th>
               <th>Godzina</th>
               <th>Opis</th>
               <th>Oddział</th>
            </thead>
            <tbody>
               <tr>
                  <td>2023-12-08</td>
                  <td>11:45:04</td>
                  <td>Przesyłka doręczona</td>
                  <td></td>
               </tr>
               <tr>
                  <td>2023-12-08</td>
                  <td>09:22:53</td>
                  <td>Przyjęcie przesyłki do punktu Pickup</td>
                  <td></td>
               </tr>
               <tr>
                  <td>2023-12-08</td>
                  <td>09:21:32</td>
                  <td>Przyjęcie przesyłki do punktu Pickup</td>
                  <td></td>
               </tr>
               <tr>
                  <td>2023-12-08</td>
                  <td>08:16:57</td>
                  <td>Wydanie do doręczenia za granicą</td>
                  <td></td>
               </tr>
               <tr>
                  <td>2023-12-08</td>
                  <td>05:24:10</td>
                  <td>Przyjęcie przesyłki w oddziale doręczenia za granicą</td>
                  <td></td>
               </tr>
               <tr>
                  <td>2023-12-07</td>
                  <td>18:35:44</td>
                  <td>Przeładunek w sortowni za granicą</td>
                  <td></td>
               </tr>
               <tr>
                  <td>2023-12-07</td>
                  <td>18:35:40</td>
                  <td>Przeładunek w sortowni za granicą</td>
                  <td></td>
               </tr>
               <tr>
                  <td>2023-12-07</td>
                  <td>01:02:02</td>
                  <td>Przekazano za granicę<br><a href="https://www.dpdgroup.com/pl/mydpd/my-parcels/search?lang=pl&amp;parcelNumber=13559300030524" class="underline" target="_blank">13559300030524</a></td>
                  <td>NCS</td>
               </tr>
               <tr>
                  <td>2023-12-06</td>
                  <td>18:10:18</td>
                  <td>Przyjęcie przesyłki w oddziale DPD </td>
                  <td>LEG</td>
               </tr>
               <tr>
                  <td>2023-12-06</td>
                  <td>15:54:08</td>
                  <td>Powiadomienie mail</td>
                  <td>WA1</td>
               </tr>
               <tr>
                  <td>2023-12-06</td>
                  <td>15:14:57</td>
                  <td>Przesyłka odebrana przez Kuriera</td>
                  <td>MDZ</td>
               </tr>
               <tr>
                  <td>2023-12-06</td>
                  <td>10:57:52</td>
                  <td>Zarejestrowano dane przesyłki, przesyłka jeszcze nienadana</td>
                  <td></td>
               </tr>
            </tbody>
         </table>
      </div>
   </fieldset>
</div>`

	xpath := `{"code":"dpd-poland","tasks":[{"type":"request","payload":"https://tracktrace.dpd.com.pl/EN/findPackage","params":{"body":"q=[track]\u0026typ=1","method":"POST","type":"xpath"},"field":{}},{"type":"query","payload":"//table[@class=\"table-track\"]/tbody/tr","field":{"path":"events","type":"array","element":{"type":"object","object":[{"path":"status","query":"./td[3]/text()[1]"},{"path":"fulldate","query":"concat(//td[1],' ', //td[2])"},{"path":"date","query":"td[1]"},{"path":"time","query":"td[2]"},{"path":"place","query":"./td[4]"},{"path":"details","query":"./td[3]/a"}]}}}]}`
	html := `{"code":"dpd-poland","tasks":[{"type":"request","payload":"https://tracktrace.dpd.com.pl/EN/findPackage","params":{"body":"q=[track]\u0026typ=1","method":"POST","type":"html"},"field":{}},{"type":"query","payload":"tbody tr","field":{"path":"events","type":"array","element":{"type":"object","object":[{"path":"status","query":"td:nth-of-type(3)"},{"path":"fulldate","query":"concat(//td[1],' ', //td[2])"},{"path":"date","query":"td:nth-of-type(1)"},{"path":"time","query":"td:nth-of-type(2)"},{"path":"place","query":"td:nth-of-type(4)"},{"path":"details","query":"td:nth-of-type(3) a"}]}}}]}`

	s1 := Scraper{}
	s2 := Scraper{}

	json.Unmarshal([]byte(xpath), &s1)
	json.Unmarshal([]byte(html), &s2)

	s1.Tasks[0].Type = "source"
	s1.Tasks[0].Payload = data

	s2.Tasks[0].Type = "source"
	s2.Tasks[0].Payload = data

	b.Run("xpath", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			args := &Args{
				ResultBuilder: NewResultBuilder(),
			}
			s1.Scrape(args)
			args.ResultBuilder.GetString()
		}
	})

	b.Run("html", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			args := &Args{
				ResultBuilder: NewResultBuilder(),
			}
			s2.Scrape(args)
			args.ResultBuilder.GetString()
		}
	})
}

func BenchmarkJsonDoc(b *testing.B) {
	data := `{"trackingsDto":null,"detailedTrackingsWithMultiMailDto":null,"detailedTrackings":[{"userTrackingItemId":null,"userTitle":null,"itemAddedDate":null,"deleteDate":null,"lastOperationViewed":false,"deleted":false,"autoAdded":false,"lastOperationViewedTimestamp":null,"payable":false,"paymentStatus":null,"paymentSystem":null,"paymentDate":null,"euv":false,"amount":null,"formless":false,"formlessData":null,"trackingItem":{"destinationCountryName":"Россия","destinationCountryNameGenitiveCase":"России","destinationCityName":"Казань","originCountryName":"Китай","originCityName":null,"mailRank":null,"mailCtg":0,"postMark":0,"insurance":null,"insuranceMoney":null,"isDestinationInInternationalTracking":true,"isOriginInInternationalTracking":true,"hiddenHistoryList":[],"futurePathList":[],"cashOnDeliveryEventsList":[],"shipmentTripInfo":{"acceptance":{"date":"2023-11-20T15:01:00.000+08:00","operationType":1,"operationAttr":0,"index":"","indexTo":"420005"},"customsPassing":null,"arrived":{"date":"2023-12-18T09:49:30.000+03:00","operationType":8,"operationAttr":2,"index":"420005","indexTo":"420005"},"hiddenInternalHistoryRecord":null,"expectedDeliveryDays":null,"startDeliveryDate":null,"expectedDeliveryDate":null,"isExpectedDeliveryDateAvailable":false},"sender":"Hong Huo Ming","recipient":"Тихонов С. А.","weight":55,"storageTime":0,"title":"Мелкий пакет из Китая","liferayWebContentId":null,"trackingHistoryItemList":[{"date":"2023-12-21T08:02:45.019+03:00","humanStatus":"Получено адресатом","humanOperationStatus":null,"operationType":2,"operationAttr":1,"countryId":643,"index":"420005","cityName":"Казань","countryName":"Россия","countryNameGenitiveCase":"России","countryCustomName":"Почта России","description":"420005, Казань","weight":55,"isInInternationalTracking":true},{"date":"2023-12-18T09:49:30.000+03:00","humanStatus":"Прибыло в место вручения","humanOperationStatus":null,"operationType":8,"operationAttr":2,"countryId":643,"index":"420005","cityName":"Казань","countryName":"Россия","countryNameGenitiveCase":"России","countryCustomName":"Почта России","description":"420005, Казань","weight":55,"isInInternationalTracking":true},{"date":"2023-11-20T15:01:00.000+08:00","humanStatus":"Принято в отделении связи","humanOperationStatus":null,"operationType":1,"operationAttr":0,"countryId":156,"index":"","cityName":null,"countryName":"Китай","countryNameGenitiveCase":"Китая","countryCustomName":"Почта Китая","description":"Китай","weight":null,"isInInternationalTracking":true}],"globalStatus":"ARCHIVED","mailType":"Мелкий пакет","mailTypeCode":5,"countryFromCode":156,"countryToCode":643,"customDuty":null,"cashOnDelivery":null,"indexFrom":"523690","indexTo":"420005","canBeOrdered":false,"canBePickedUp":false,"deliveryOrderDate":null,"commonStatus":"Вручено 21 декабря","firstOperationDate":1700463660000,"lastOperationDate":1703134965019,"lastOperationTimezoneOffset":10800000,"lastOperationDateTime":"2023-12-21T08:02:45.019+00:00","acceptanceOperationDateTime":"2023-11-20T07:01:00.000+00:00","complexCode":null,"complexType":null,"complexDeliveryMethod":null,"barcode":"ZA298249753CN","barcodeImage":null,"notificationBarcode":null,"notificationTitle":null,"sourceBarcode":null,"sourceTitle":null,"endStorageDate":null,"lastDayInOps":null,"hasBeenGiven":true,"lastOperationAttr":1,"lastOperationType":2,"id":null,"postmarkText":"","mailTypeText":"Мелкий пакет","customsPaymentStatus":null,"customsOperatorDuty":null,"completeness":false,"tpoCustomPayment":null,"customsOperatorDutyFee":null,"customsOperatorDutyFeeOnline":null,"customsPayment":null,"customsPayments":[],"mailCtgText":"Простое","mailRankText":null,"returnRate":null,"redispatchRate":null},"alienOrder":null,"opmInfo":null,"paymentSchema":"","officeSummary":null,"postmanDeliveryInfo":null,"courierDeliveryInfo":null,"omsDeliveryInfo":null,"omsPickupInfo":null,"pickupAvailableAddress":null,"formF22Params":{"MailTypeText":"Мелкий пакет","senderAddress":"Китай","RecipientIndex":"420005","WeightGr":55,"SmallPackage":true,"SendingType":"SmallPackage","SimplePackage":true,"MailCtgText":"Простое","state":null,"PostId":"ZA298249753CN"},"attorneyServiceIsEnabled":true,"deliveryAvailable":false,"timeSlotInfo":null,"delayInsuranceInfo":{"paymentStatus":"UNAVAILABLE","productInfo":null,"descriptionUrl":null},"refusalInsuranceInfo":{"paymentStatus":"UNAVAILABLE","productInfo":null,"descriptionUrl":null},"officeReservationUrl":"","parcelStatus":null,"euvStatus":"missing","opsAppointmentAvailable":false,"isMailCopyAvailable":false}]}`

	xpath := `{"code":"code","tasks":[{"type":"request","payload":"http://www.pochta.ru/api/tracking/api/v1/trackings/by-barcodes?language=ru\u0026track-numbers=[track]","params":{"type":"json.xpath"},"field":{}},{"type":"query","payload":"//detailedTrackings//recipient","field":{"path":"Recipient"}},{"type":"query","payload":"//detailedTrackings//sender","field":{"path":"Sender"}},{"type":"query","payload":"//detailedTrackings//trackingHistoryItemList/*","field":{"path":"events1"}},{"type":"query","payload":"//detailedTrackings//trackingHistoryItemList/*","field":{"path":"events2","type":"array"}},{"type":"query","payload":"//detailedTrackings//trackingHistoryItemList/*","field":{"path":"events3","type":"array","element":{"type":"object","object":[{"path":"status","query":"./humanStatus"},{"path":"date","query":"./date"},{"path":"place","query":"./cityName"},{"path":"zip","query":"./index"},{"path":"meta.description","query":"./description"},{"path":"meta.service","query":"./countryCustomName"}]}}},{"type":"query","payload":"//detailedTrackings//trackingHistoryItemList/*[2]","field":{"path":"events4","type":"object","object":[{"path":"status","query":"./humanStatus"},{"path":"date","query":"./date"},{"path":"place","query":"./cityName"},{"path":"zip","query":"./index"},{"path":"meta.description","query":"./description"},{"path":"meta.service","query":"./countryCustomName"}]}}]}`
	gjson := `{"code":"code","tasks":[{"type":"request","payload":"http://www.pochta.ru/api/tracking/api/v1/trackings/by-barcodes?language=ru\u0026track-numbers=[track]","params":{"type":"json"},"field":{}},{"type":"query","payload":"detailedTrackings.0.trackingItem.recipient","field":{"path":"Recipient"}},{"type":"query","payload":"detailedTrackings.0.trackingItem.sender","field":{"path":"Sender"}},{"type":"query","payload":"detailedTrackings.0.trackingItem.trackingHistoryItemList","field":{"path":"events1"}},{"type":"query","payload":"detailedTrackings.0.trackingItem.trackingHistoryItemList","field":{"path":"events2","type":"array"}},{"type":"query","payload":"detailedTrackings.0.trackingItem.trackingHistoryItemList","field":{"path":"events3","type":"array","element":{"type":"object","object":[{"path":"status","query":"humanStatus"},{"path":"date","query":"date"},{"path":"place","query":"cityName"},{"path":"zip","query":"index"},{"path":"meta.description","query":"description"},{"path":"meta.service","query":"countryCustomName"}]}}},{"type":"query","payload":"detailedTrackings.0.trackingItem.trackingHistoryItemList.1","field":{"path":"events4","type":"object","object":[{"path":"status","query":"humanStatus"},{"path":"date","query":"date"},{"path":"place","query":"cityName"},{"path":"zip","query":"index"},{"path":"meta.description","query":"description"},{"path":"meta.service","query":"countryCustomName"}]}}]}`

	s1 := Scraper{}
	s2 := Scraper{}

	json.Unmarshal([]byte(xpath), &s1)
	json.Unmarshal([]byte(gjson), &s2)

	s1.Tasks[0].Type = "source"
	s1.Tasks[0].Payload = data

	s2.Tasks[0].Type = "source"
	s2.Tasks[0].Payload = data

	b.Run("xpath", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			args := &Args{
				ResultBuilder: NewResultBuilder(),
			}
			s1.Scrape(args)
			args.ResultBuilder.GetString()
		}
	})

	b.Run("gjson", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			args := &Args{
				ResultBuilder: NewResultBuilder(),
			}
			s2.Scrape(args)
			args.ResultBuilder.GetString()
		}
	})
}
