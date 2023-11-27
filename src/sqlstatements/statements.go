package sqlstatements

// MoveTablesSQL contains the SQL statement to move tables to the user's schema
const MoveTablesSQL = `
DO
$$
DECLARE
row record;
BEGIN
FOR row IN SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname = 'public' AND tableowner = current_user
LOOP
RAISE NOTICE 'Processing table: %', row.tablename;
EXECUTE 'ALTER TABLE public.' || quote_ident(row.tablename) || ' SET SCHEMA ' || current_user || '';
END LOOP;
END;
$$;
`

// MoveSequencesSQL is SQL which moves sequences to the user's schema
const MoveSequencesSQL = `
DO
$$
DECLARE
row record;
BEGIN
FOR row IN SELECT sequence_name FROM information_schema.sequences WHERE sequence_schema = 'public' AND sequence_catalog IN (current_user, current_database())
LOOP
RAISE NOTICE 'Processing sequence: %', row.sequence_name;
EXECUTE 'ALTER SEQUENCE public.' || quote_ident(row.sequence_name) || ' SET SCHEMA ' || current_user || '';
END LOOP;
END;
$$;
`

// MoveViewsSQL is SQL which moves views to the user's schema
const MoveViewsSQL = `
DO
$$
DECLARE
row record;
BEGIN
FOR row IN SELECT viewname FROM pg_catalog.pg_views WHERE schemaname = 'public' AND viewowner = current_user
LOOP
RAISE NOTICE 'Processing view: %', row.viewname;
EXECUTE 'ALTER VIEW public.' || quote_ident(row.viewname) || ' SET SCHEMA ' || current_user || '';
END LOOP;
END;
$$;
`
