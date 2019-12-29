import React, { Component } from 'react';
import './App.css';
import 'leaflet/dist/leaflet.css';
import { Map, Marker, Popup, TileLayer } from 'react-leaflet';
import axios from 'axios';


class IncidentMap extends Component {
  constructor(props) {
    super(props);

    this.state = {
      incident: {
        address: {},
        apparatus: [],
        description: {},
        fire_department: {},
      },
    };
  }

  componentDidMount() {
    this.cancelToken = axios.CancelToken.source();

    this.loadIncident();
  }

  async loadIncident() {
    try {
    const response = await axios.get('/incident', { cancelToken: this.cancelToken.token });
      this.setState({ incident: response.data });
    } catch (error) {
      console.error(error);
    }
  }

  render() {

    const incident = { ...this.state.incident };

    const position = [incident.address.latitude || 39, incident.address.longitude || -98];

    const mapMarkers = incident.apparatus.map(a =>
      (Object.keys(a.unit_status)).forEach(u =>
      <Marker
        position={[a[u].latitude, a[u].longitude]}
      >
        <Popup>
          {a.unit_type}<br/>
          {u}<br/>
          {a[u].timestamp}
        </Popup>
      </Marker>
      )
    );

    return (
    <div className="App">
      <header className="App-header">
        <p>{incident.fire_department ? incident.fire_department.name : ''}</p>
      </header>
      <p>{incident.description ? incident.description.comments : ''}</p>
      <div style={{ height: '500px', width: '100%' }}>
        <Map center={position} zoom={13}>
          <TileLayer
             attribution='&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'
             url='https://{s}.tile.osm.org/{z}/{x}/{y}.png'
           />
        {mapMarkers}
        </Map>
      </div>
    </div>
    );
  }
}

export default IncidentMap;
