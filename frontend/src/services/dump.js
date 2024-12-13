import axios from 'axios';
const baseUrl = "http://localhost:8083";

const getDump = async () => {
  const startDate = new Date();
  startDate.setDate(startDate.getDate() - startDate.getDay() + 1);
  const endDate = new Date(startDate);
  endDate.setDate(endDate.getDate() + 6);
  const formatDate = (date) => date.toISOString().split('T')[0];
  const formattedStartDate = formatDate(startDate);
  const formattedEndDate = formatDate(endDate);
  try {
    const response = await axios.get(`${baseUrl}/summary/${formattedStartDate}/${formattedEndDate}`, {
      responseType: 'blob',
    });
    const returnedData = URL.createObjectURL(response.data)
    return returnedData;
  } catch (error) {
    console.error("Error fetching dumps:", error);
    throw error;
  }
}

export { getDump }