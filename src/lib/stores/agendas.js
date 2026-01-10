import {writable} from "svelte/store";
import axios from "axios";

export const agendas = writable([])
export const agendasError = writable("")

const apiBaseUrl = "/config_api/agendas"

/*
***** DATA FORMAT EXAMPLE
* {
*  "id": "6ce750e5-05b8-4c4c-8fee-d9e381dbf364",
*  "ucaId": 56529,
*  "name": "M1 - Tutorat L2"
* }
*/

export function getAgendas() {
    agendasError.set("")
    axios.get(`${apiBaseUrl}`)
        .then((res) => {
            agendas.set(res.data)
        })
        .catch((err) => {
            console.log("An error has occurred while retrieving agendas")
            console.log(err)
            if(err.response?.data?.message){
                agendasError.set(JSON.stringify(err.response.data.message))
            } else {
                agendasError.set(JSON.stringify(err))
            }
        })
}

export function getAgenda(id) {
    agendasError.set("")
    return axios.get(`${apiBaseUrl}/${id}`)
        .then((res) => {
            return Promise.resolve(res.data)
        })
        .catch((err) => {
            console.log("An error has occurred while retrieving agenda")
            console.log(err)
            if(err.response?.data?.message){
                agendasError.set(JSON.stringify(err.response.data.message))
            } else {
                agendasError.set(JSON.stringify(err))
            }
            return Promise.reject(err)
        })
}

export function postAgendas(ucaId, name) {
    agendasError.set("")
    let parsed = parseInt(ucaId)
    return axios.post(`${apiBaseUrl}`,{
        ucaId: isNaN(parsed) ? ucaId : parsed,
        name: name
    })
        .then((res) => {
            getAgendas()
            return Promise.resolve()
        })
        .catch((err) => {
            console.log("An error as occurred while posting agenda")
            console.log(err)
            if(err.response?.data?.message){
                agendasError.set(JSON.stringify(err.response.data.message))
            } else {
                agendasError.set(JSON.stringify(err))
            }
            return Promise.reject(err)
        })
}

export function putAgenda(id, ucaId, name) {
    agendasError.set("")
    let parsed = parseInt(ucaId)
    return axios.put(`${apiBaseUrl}/${id}`,{
        ucaId: isNaN(parsed) ? ucaId : parsed,
        name: name
    })
        .then((res) => {
            getAgendas()
            return Promise.resolve()
        })
        .catch((err) => {
            console.log("An error as occurred while putting agenda")
            console.log(err)
            if(err.response?.data?.message){
                agendasError.set(JSON.stringify(err.response.data.message))
            } else {
                agendasError.set(JSON.stringify(err))
            }
            return Promise.reject(err)
        })
}

export function deleteAgenda(id) {
    agendasError.set("")
    return axios.delete(`${apiBaseUrl}/${id}`)
        .then((res) => {
            getAgendas()
            return Promise.resolve()
        })
        .catch((err) => {
            console.log("An error as occurred while deleting agenda")
            console.log(err)
            if(err.response?.data?.message){
                agendasError.set(JSON.stringify(err.response.data.message))
            } else {
                agendasError.set(JSON.stringify(err))
            }
            return Promise.reject(err)
        })
}