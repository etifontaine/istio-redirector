import { useState } from 'react'
import getConfig from 'next/config'

const { publicRuntimeConfig } = getConfig()

export function useGetLayout() {

	const [redirectionType, setRedirectionType] = useState("")
	const [isSuccess, showSuccess] = useState(false)
	const [CSVFile, setCSVFile] = useState()
	const [formData, setFormData] = useState({})
	const [alert, setAlert] = useState({ isVisible: false })
	const [CSVData, setCSVData] = useState([])

	const clearData = () => {
		setCSVFile()
		setCSVData([])
	}

	const handleChangeRedirectionType = (e) => {
		setRedirectionType(e.target.value);
	}

	const clearAlert = () => setAlert({ isVisible: false, content: null })

	const sendData = () => {
		if (!formData.redirection_name) {
			setAlert({ isVisible: true, content: 'Redirection name must be filled' })
			return
		}
		const formDataValues = new FormData()
		formDataValues.append('csv_file', CSVFile);
		formDataValues.append('redirection_name', formData.redirection_name)
		formDataValues.append('redirection_type', redirectionType)
		fetch(`${publicRuntimeConfig.API_URL}/api/csv/upload`, { method: 'POST', body: formDataValues })
			.then(async response => {
				const file = await response.blob()
				const href = window.URL.createObjectURL(file);
				const link = document.createElement('a');
				link.href = href;
				link.setAttribute('download', 'config.yaml'); //or any other extension
				document.body.appendChild(link);
				link.click();

				clearData()
			})
			.then(() => showSuccess(true))
	}

	return {
		setCSVFile,
		setCSVData,
		CSVFile,
		CSVData,
		clearData,
		sendData,
		handleChangeRedirectionType,
		redirectionType,
		isSuccess,
		showSuccess,
		formData,
		setFormData,
		alert,
		setAlert,
		clearAlert
	}
}
