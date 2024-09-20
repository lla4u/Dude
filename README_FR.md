# Dude - Explorateur de données utilisateur Dynon

> Chaque écran SkyView HDX peut agir comme un écran de vol primaire (PFD) avec vision synthétique, un
Système de surveillance du moteur (EMS) et carte mobile dans une variété d'écrans personnalisables
mises en page. Les données proviennent de divers modules et appareils connectés.
>
> SkyView HDX affiche, enregistre et stocke les informations de vol dans plusieurs journaux de données qui peuvent être
exporté pour analyse par le propriétaire, et un journal de données haute résolution qui peut être utilisé par Dynon
pour le dépannage. 
>
> Cet outil offre un moyen simple et efficace d'éclairer les datalogs Dynon pour :  
> - Fournir un historique de vol à long terme,
> - Fournir la capacité d'afficher la carte du vol et les paramètres associés pour améliorer l'utilisation et la sécurité du pilote,
> - Fournir des informations précises comme la vitesse moyenne à envisager pour préparer le vol, la vitesse moyenne d'atterrissage,
> - ...

## Quelles sont les dépendances
> Cet outil utilise :
> - Docker pour les conteneurs et la gestion du réseau
> - Base de données Influxdb pour un stockage chronodaté à long terme
> - Gafana pour la présentation des données
> - Un shell Unix et un programme dédié (go) pour l'intégration des journaux de données Dynon dans la base de données.
>
> Tous les éléments sélectionnés sont open source et gratuits pour un usage personnel.

## Ce qui est proposé
> Une interface web locale pour interroger les données de vol telles que :
> ![Capture d'écran de l'interface Web.](https://github.com/lla4u/Dude-Influx-Grafana/blob/main/Screenshots/Screenshot_web_Interface.png)
>
> ![Capture d'écran de l'interface web 2.](https://github.com/lla4u/Dude-Influx-Grafana/blob/main/Screenshots/Screenshot_web_interface_2.png)

> [!NOTE]
> Il est possible de construire une présentation personnalisée.

# Procédure d'installation

## 1 - Installation de Docker (si ce n'est pas encore fait)
En fonction de votre système d'exploitation, téléchargez et installez le logiciel Docker.
https://docs.docker.com/get-docker/  

Vidéo : https://www.youtube.com/watch?v=mS26N5cLBe8&ab_channel=CodersArcade  


## 2 - Construire la stack Dude
```
1. Où sera placée la stack :
   Ouvrez un terminal, cmd et créez votre répertoire d'installation personnel 
     cd /home/lla 
     mkdir dude 
   puis passez à : 
     cd dude

2. Cloner le dépôt github :
   git clone https://github.com/lla4u/Dude-Influx-Grafana.git
   ou
   Téléchargez et décompressez l'archive zip téléchargée depuis github.

3. Changer de répertoire pour Dude-Influx-Grafana
   cd Dude-Influx-Grafana

4. Construisez la stack Docker à l'aide du terminal : 
   docker-compose --env-file config.env up --build -d 
   ou (pour la version récente de Docker) 
   docker compose --env-file config.env up --build -d 

   après un certain temps (principalement en fonction de la bande passante de votre réseau), 3 conteneurs seront créés et disponibles.

5. Vérifiez :
   exécuter docker ps depuis le terminal

   Ayant 3 conteneurs en marche, vous êtes prêt à aller plus loin...
```
> [!IMPORTANT]
> Les données persistantes (InfluxDB & Grafana) seront stockées dans un répertoire local (Docker). La suppression de celui-ci entraînera une perte de données.

## 2 - Démarrage / Arrêt de la stack Dude
> Le démarrage ou l'arrêt de la pile Dude peut être réalisé en utilisant :
> 1. Tableau de bord Docker
> - Démarrage
> ![Capture d'écran du tableau de bord Docker démarrant.](https://github.com/lla4u/Dude-Influx-Grafana/blob/main/Screenshots/Screenshot_docker_dashboard_start.png)
>
> - Arrêt
> ![Capture d'écran du tableau de bord Docker arrêté.](https://github.com/lla4u/Dude-Influx-Grafana/blob/main/Screenshots/Screenshot_docker_dashboard_stop.png)
>
> 2. Ligne de commande
> - Ouvrez le terminal et accédez au répertoire Dude-Influx-Grafana
> - cd /home/lla/dude/Dude-Influx-Grafana
>
> - Démarrage
> - Docker compose up
>
> - Arrêt
> - Docker compose down

# Ajout de datalogs HDX à la solution
> L'ajout de journaux de données dans l'outil est un processus en deux étapes :
> - Tout d'abord, collectez le journal de données du HDX
> - Deuxièmement, importez les journaux de données dans Influxdb à l'aide du conteneur dude-cli.

## Collecte du journal de données du HDX
> La collecte des données du HDX est assez triviale et nécessite une clé USB branchée sur le Dynon :
> (J'utilise la même clé usb que pour les plaques et mises à jour des cartes...)
> 1. Allumez votre HDX
> 
> 2. Appuyez simultanément sur les boutons 7 et 8 pendant quelques secondes pour démarrer l'écran de configuration Dynon.
> 
> 3. Accédez à SYSTEM SOFTWARE -> EXPORT USER DATA LOGS 
> 
> 4. (Facultatif) Définir l'étiquette
> 
> 5. Exporter en appuyant sur le bouton 8
> 
> Vidéo : https://www.youtube.com/watch?v=fS6H_8gNd90&ab_channel=RobertHamilton

> [!IMPORTANT]
> Le stockage des journaux de données Dynon est limité et réécrit au fil du temps. Collectez donc des datalogs toutes les 8 heures de vol environ ou acceptez de perdre des informations.

## Importation du journal de données dans la pile Dude
> 1. Copiez le(s) fichier(s) csv de la clé USB (USER_DATA_LOG.csv) dans le répertoire Datalogs
> 
> 2. Depuis Terminal ou cmd ou PowerShell (windows), exécutez : 
>    - docker exec -it mec-cli /bin/ash
> 
> 3. Exécutez :
>    - ./dude-cli 
>   Mode verbeux facultatif :
>    - ./dude-cli  -v
> 
> ![Capture d'écran de dude-cli.](https://github.com/lla4u/Dude-Influx-Grafana/blob/main/Screenshots/Screenshot_dude-cli.png)
> Remarque sur la capture d'écran :  
> - L'intégration des journaux de données a nécessité 41.24 secondes (principalement en raison de mauvaises écritures synchrones d'Influxdb)  
> - Le fichier soumis contient 164897 lignes CSV  
> - L'importation a conservé 39727 lignes enregistrées dans la base de données.

> [!NOTE]
> dude-cli importe uniquement les nouveau datalogs !. Le fichier imported.txt est la reference des datalogs déja traités.

# Feuille de route
- [x] Correction d'une erreur d'écriture d'InfluxDB ayant un nombre de fichiers énorme lors de l'importation. (Corrigé en utilisant l'écriture synchrone et allez à la place de nodejs)
- [x] Importer automatiquement la configuration des tableaux de bord Grafana et de la source de données 
- [ ] Utilisez la variable Grafana pour vous aider à trouver les vols enregistrés dans InfluxDB.
- [x] Améliorer les performances de Dude-cli.
- [x] Améliorer l'interface utilisateur de Dude-cli.
- [ ] Créer un document pour aider les utilisateurs à utiliser l'interface utilisateur Grafana et consulter les journaux de données Dynon.

Ayez des vols en toute sécurité.

Laurent