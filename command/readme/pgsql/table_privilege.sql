select n.nspname, c.relname,
       has_table_privilege('username', c.oid, 'SELECT') as can_select,
       has_table_privilege('username', c.oid, 'INSERT') as can_insert,
       has_table_privilege('username', c.oid, 'UPDATE') as can_update,
       has_table_privilege('username', c.oid, 'DELETE') as can_delete
  from pg_class c
  join pg_namespace n on c.relnamespace = n.oid
where c.relkind = 'r'
  and n.nspname = 'public'
  and c.relname = 'table_name';