- [x] URL short (changes DB structure, adds Auto_incr ID used to encode
       and decode into base62) - this makes UUID redundant (on low
       scale) -> `Solved via Snowflake like ID Generation`
       
- [x] Migrations ? -> `DB structure wont be changing, init script will do`
       
- [ ] Cache (Recently pushed pastes cache for couple hours)?
       
- [ ] Docker file for Backend on same network in compose
       
- [x] Few MB limit to pastes -> `Handling on http handler boundary`

- [x] Create Cleaup cron job / Garbage Collector

- [ ] FrontEnd (Ugh)
