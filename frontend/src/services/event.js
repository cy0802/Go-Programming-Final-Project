import axios from 'axios';
const baseUrl = "http://localhost:8083";

const getEvents = async () => {
  try {
    const response = await axios.get(`${baseUrl}/diary`);
    var returnedData = [];
    response.data.forEach((element) => {
      element.events.forEach((event) => {
        returnedData.push({
          id: event.id,
          date: element.date,
          title: event.title,
          diary: event.content,
        });
      });
    });
    return returnedData;
  } catch (error) {
    console.error("Error fetching events:", error);
    throw error;
  }
}

const createEvent = async (event) => {
  try {
    const date = event.date;
    const response = await axios.post(`${baseUrl}/diary/${date}`, {
      title: event.title,
      content: event.diary,
    });
    return response.data;
  } catch (error) {
    console.error("Error creating event:", error);
    throw error;
  }
}

const updateEvent = async (event) => {
  try {
    const id = event.id;
    const response = await axios.put(`${baseUrl}/diary/${id}`, {
      id: event.id,
      title: event.title,
      content: event.diary,
    });
    return response.data;
  } catch (error) {
    console.error("Error updating event:", error);
    throw error;
  }
}

const deleteEvent = async (id) => {
  try {
    const response = await axios.delete(`${baseUrl}/diary/${id}`);
    return response.data;
  } catch (error) {
    console.error("Error deleting event:", error);
    throw error;
  }
}

export {
  getEvents,
  createEvent,
  updateEvent,
  deleteEvent,
}