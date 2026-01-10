import {writable, get} from "svelte/store";
import axios from "axios";

export const events = writable([])
export const eventsError = writable("")

const apiBaseUrl = "/timetable_api/events"

/*
***** DATA FORMAT EXAMPLE
* {
*  "id": "62a2beca-26cf-45bf-aa82-4cf5b14922fd",
*  "agendaIds": [
*     "d5c60e7a-10cd-4aec-9ea5-96d071ba824b"
*  ],
*  "uid": "ADE60323032342d323032352d5543412d36303334342d302d32",
*  "description": "\\n\\nM1 GROUPE 1 langue\\nPAILLOUX MARIE\\n\\n(Updated :26/11/202 4 09:51)",
*  "name": "TD Entrepôt de données - G1",
*  "start": "2025-01-23T15:45:00+01:00",
*  "end": "2025-01-23T17:45:00+01:00",
*  "location": "IS_A104",
*  "lastUpdate": "2024-11-26T09:51:00+01:00"
* }
*/

export function getEvents() {
    eventsError.set("")
    axios.get(`${apiBaseUrl}`)
        .then((res) => {
            events.set(res.data)
        })
        .catch((err) => {
            console.log("An error has occurred while retrieving events")
            console.log(err)
            if(err.response?.data?.message){
                eventsError.set(JSON.stringify(err.response.data.message))
            } else {
                eventsError.set(JSON.stringify(err))
            }
        })
}