# Uzantoj #

El objetivo de la aplicación es la creación y borrado usuarios en LDAP. Se asume que estos son usuarios de Mail (Postfix, Courier, Dovecot, etc).

Básicamente se puede:

* Crear y borrar dominios.
* Crear y borrar cuentas (usuarios de LDAP y Maildirs).

En este directorio están:

* **sql:** el directorio con los esquemas de BD. Usar la más nueva (uzantoj-v2.sql).
* **uzantoj:** el directorio con el proyecto Revel.

---

## Funcionalidad ##

* Dominios: son los dominios de email.
* Usuarios: son los usuarios del sistema (los que están en el MySQL de este sistema, no los guardados en el LDAP).
* Cola de trabajos: lo que decíamos de usar RabbitMQ, pero creo que no es necesario.
* Mi perfil: es la edición del usuario que está logueado.
* Ayuda: debería ser una página estática con algunos tips que consideremos necesarios.
* Contacto: lo pensé como un form de contacto.

Para un dominio podemos:

* Verlo
* Borrarlo
* Editarlo no, porque no hay nada que editar
* Listar los usuarios de email (los que están guardados en el LDAP)
* Agregar usuarios de email (los que se guardan en el LDAP)
* Vincular administrador al dominio. Acá se vincula un usuario del sistema con un dominio. Este vínculo existe sólo en el MySQL y no se refleja en el LDAP. Es la tabla domains_users.
* Desvincular administrador. Es lo contrario al punto anterior.

---

## Requerimientos ##

La aplicación está escrita en GoLang usando el _framework_ Revel. Por lo tanto, para compilarla es necesario GoLang. Para correrla, es necesario Revel.

Además, se asume que la base de datos está en un servidor de MySQL. Y que los usuarios del dominio están en un OpenLDAP. Por último, para la creación y borrado de directorios (como Maildirs), la aplicación se conecta por medio de SSH al servidor y crea los directorios.

---

## Configuración del entorno ##

Hay que configurar 3 cosas: base de datos MySQL, servidor de LDAP y acceso por SSH a la máquina donde van a estar los Maildirs.

### Base de datos ###

Se usa para guardar los usuarios del sistema (de la aplicación). Los usuarios de email se guardan en el LDAP. 

Crear la base de datos y carga el dump _sql/uzantoj-v2.sql_. Luego, en el script _run.sh_ escribir los datos de acceso:

	export UZANTOJ_DB_NAME="db-name"
	export UZANTOJ_DB_USER="usuario"
	export UZANTOJ_DB_PASSWORD="password"
	export UZANTOJ_DB_SERVER="servidor.local"
	export UZANTOJ_DB_PORT="3306"

### Servidor de LDAP ###

Aquí se guardan los dominios y los usuarios de mail de esos dominios. Los _object classes_ usados son:

Para los dominios ("dc=testing.com.ar,dc=mydomain"): 

* _posixAccount_, 
* _organization_, 
* _top_ y
* _dcObject_

Para el grupo de usuarios ("ou=users,dc=testing.com.ar,dc=mydomain"): 

* _organizationalUnit_ y
* _top_

Para los usuarios ("uid=myuser_testing_com_ar,ou=users,dc=testing.com.ar,dc=mydomain"): 

* _posixAccount_ y 
* _inetOrgPerson_

Para el caso de Debian, los _schemas_ que vienen con OpenLDAP ya las implementan, por lo que no es necesario realizar ningún cambio.

Las variables a configurar en _run.sh_ son:

	UZANTOJ_LDAP_SERVER="servicios-desarrollo.local"
	UZANTOJ_LDAP_PORT="389"
	UZANTOJ_LDAP_BIND="cn=admin,dc=mydomain"
	UZANTOJ_LDAP_PASSWORD="my-password"
	UZANTOJ_LDAP_ROOT="dc=nodomain"


### Acceso por SSH ###

Se requiere acceso por SSH para crear los directorios para los dominios y los Maildir de los usuarios. En este caso tendremos:

* **Mailserver:** es el servidor donde se van a crear los directorios.
* **AppServer:** es el servidor donde corre la aplicación.
* **uzantoj:** es el usuario que se va a crear en Mailserver. La conexión es por SSH desde AppServer a Mailserver, usando un par de claves pública/privada y sin que pida contraseña.

Los pasos descritos a continuación son genéricos. Una mejora es restringir los comandos que puede ejecutar el usuario utilizando _sudo_.

Las variables en _run.sh_ son:

	UZANTOJ_MAILSERVER_HOSTNAME="servicios-desarrollo.local"
	UZANTOJ_MAILSERVER_USER="uzantoj"
	UZANTOJ_MAILSERVER_HOMEDIR="/var/vdomains"


#### Crear el usuario en Mailserver ####

Creamos el usuario:

	adduser uzantoj

#### Crear el directorio para los dominios y Maildirs ####

Creamos el directorio principal:

	mkdir -p /var/vdomains

Lo hacemos propiedad de _uzantoj_:

	chown uzantoj:uzantoj /var/vdomains/

Y acomodamos los permisos:

	chmod 750 /var/vdomains/

#### Copiamos la clave púbica ####

Copiamos la clave pública a _.ssh/authorized_keys_.

Si el directorio _.ssh_ no existe, lo podemos crear (debe estar en el $HOME del usuario _uzantoj_):

	mkdir .ssh
	chmod 700 .ssh

La configuración de SSH que viene con Debian permite el login sin contraseña uando se utiliza la clave pública. De lo contrario, hay que configurarlo.

---

## Correr con run.sh ##

Copiar el script _run.sh.example_ de muestra a _run.sh_. Configurar con los valores correctos según el ambiente y ejecutarlo.

	cp run.sh.example run.sh

La configuración de la aplicación (_conf/app.conf_)  espera ciertas variables de entorno (como _$UZANTOJ_DB_NAME_), por lo que es más fácil setear el valor de estas variables en el script _run.sh_ y correr la aplicación desde ahí:

	./run.sh

En este caso se asume que ya se tiene un entorno configurado (base de datos MySQL, servidor LDAP y acceso por SSH a un servidor donde crear los directorios).

---

## Correr con Docker ##

Otra opción es utilizar Docker para levantar la aplicación en un ambiente de testing. Notar que todos los datos se perderán cuando se bajen los contenedores.

---

## Instalación ##

Independientemente del método usado para ejecutar la aplicación, luego hay que ingresar a *app-ip:app-port/install* (por ejemeplo, *localhost:9000/install*). Hay varias cosas para completar:

* Nombre de usuario: puede ser cualquier. Es el "admin" inicial.
* Contraseña: deben coincidir las dos contraseñas.
* Número de GID: es el group id a partir del cual se van a crear los dominios.
* Número de UID: es el user id a partir del cual se van a crear usuarios.

Por ejemplo, si GID es 1234, entonces el dominio Dom1 va a tener el group id 1234, y el Dom2 el 1235, etc. Y con los usuarios lo mismo. 

Al final nos lleva al login.

