CREATE SCHEMA incident_report;

CREATE TABLE IF NOT EXISTS incident_report.incident_data (
  incident_number text,
  address json,
  apparatus json,
  description json,
  fire_department json,
  version text,
  uploaded timestamp,
  PRIMARY KEY (incident_number)
);

CREATE OR REPLACE FUNCTION incident_report.get_last_incident() RETURNS TABLE(address json, apparatus json, description json, fire_department json, version text)
  LANGUAGE sql
  AS $$
  SELECT (address, apparatus, description, fire_department, version) FROM incident_report.incident_data
  ORDER BY uploaded DESC
  LIMIT 1;
$$;

CREATE OR REPLACE FUNCTION incident_report.upload_report(incident_number text, address json, apparatus json, description json, fire_department json, version text) RETURNS void
  LANGUAGE plpgsql
  AS $$
  DECLARE
    time_now timestamp;
  BEGIN
    time_now := now();
  INSERT INTO incident_report.incident_data(incident_number, address, apparatus, description, fire_department, version, uploaded)
  VALUES ($1, $2, $3, $4, $5, $6, time_now)
  /*
   *ON CONFLICT ON CONSTRAINT incident_data_pkey
   *DO UPDATE SET
   *  address = EXCLUDED.address,
   *  apparatus = EXCLUDED.address,
   *  fire_department = EXCLUDED.fire_department,
   *  version = EXCLUDED.version
   */
    ;
  END
$$;
