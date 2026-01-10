import {writable} from "svelte/store";
import axios from "axios";

export const alerts = writable([])
export const alertsError = writable("")

const apiBaseUrl = "/config_api/alerts"

/*
***** DATA FORMAT EXAMPLE
* {
*   "id": "3aad8a52-2777-40b9-a641-51ede4d55bdc",
*   "email": "email@etu.uca.fr",
*   "agendaId": "6ce750e5-05b8-4c4c-8fee-d9e381dbf364"
* }
*/

export function getAlerts() {
    alertsError.set("")
    axios.get(`${apiBaseUrl}`)
        .then((res) => {
            alerts.set(res.data)
        })
        .catch((err) => {
            console.log("An error has occurred while retrieving alerts")
            console.log(err)
            if(err.response?.data?.message){
                alertsError.set(JSON.stringify(err.response.data.message))
            } else {
                alertsError.set(JSON.stringify(err))
            }
        })
}

export function postAlerts(email, agendaId) {
    alertsError.set("")
    return axios.post(`${apiBaseUrl}`,{
        email: email,
        agendaId: agendaId
    })
        .then((res) => {
            getAlerts()
            return Promise.resolve()
        })
        .catch((err) => {
            console.log("An error as occurred while posting alert")
            console.log(err)
            if(err.response?.data?.message){
                alertsError.set(JSON.stringify(err.response.data.message))
            } else {
                alertsError.set(JSON.stringify(err))
            }
            return Promise.reject(err)
        })
}

export function putAlert(id, email, agendaId) {
    alertsError.set("")
    return axios.put(`${apiBaseUrl}/${id}`,{
        email: email,
        agendaId: agendaId
    })
        .then((res) => {
            getAlerts()
            return Promise.resolve()
        })
        .catch((err) => {
            console.log("An error as occurred while putting alert")
            console.log(err)
            if(err.response?.data?.message){
                alertsError.set(JSON.stringify(err.response.data.message))
            } else {
                alertsError.set(JSON.stringify(err))
            }
            return Promise.reject(err)
        })
}

export function deleteAlert(id) {
    alertsError.set("")
    return axios.delete(`${apiBaseUrl}/${id}`)
        .then((res) => {
            getAlerts()
            return Promise.resolve()
        })
        .catch((err) => {
            console.log("An error as occurred while deleting alert")
            console.log(err)
            if(err.response?.data?.message){
                alertsError.set(JSON.stringify(err.response.data.message))
            } else {
                alertsError.set(JSON.stringify(err))
            }
            return Promise.reject(err)
        })
}