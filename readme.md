# PRR Labo 2

## Étudiants
DELAY Jérémy _(elvildor)_  


## Contenu
Rien ne marche...  

Je n'ai pas réussi à faire fonctionner mon environnement de manière durable. À chaque fois que j'arrivais à le faire fonctionner, je me reçevais un erreur du type _```malformed module path "config": missing dot in first path element```_ qui est des plus floues et qui comporte de multpiles cause selon mes diverses recherches.  

J'ai essayé d'implémenter la base de l'algorithme de Lampaort avec des méthode qui pourrait lancer les messages REQ, ACK et REL si j'avais réussi à mettre en place les communications TCP.  

J'ai aussi mis en place un début de fichier de configuration, mais je n'ai bien entendu jamais pu le tester donc j'ignore si j'si eu la bonne approche ou non...


## Divers
Je n'ai rien compris à ce labo et le temps supplémentaire ne m'as même pas plus aidé que cela comtrairement à ce que je pensais. C'était un bien mauvais pari de ma part. Je ferai tout ce que je pourrai pour me rattraper avec le labo suivant vu que j'aurais la possibilité de repartir d'une base saine.  

Désolé de vous avoir fait perdre votre temps avec la correction de ce déchet que je rends comme laboratoire.


## Auto-évaluation
| # |pts max|pts perso|étudiant|  
|:-:|:-----:|:-------:|:-------|  
| 1 |   5   |    1    | Delay  |  
| 2 |  10   |    0    | Delay  |  
| 3 |  15   |    3    | Delay  |  
| 4 |  10   |    2    | Delay  |  
| 5 |   5   |    3    | Delay  |  
| 6 |   5   |    0    | Delay  |  


# PRR Labo 1

## Étudiants
DELAY Jérémy _(elvildor)_  
MAYO CARTES Adrian _(Haxos)_  

## Mode d'emploi
### Lancer le serveur
Pour lancer le serveur en ligne de commande, il faut se placer dans le dossier contenant le fichier **server.go**.  
Ensuite, il faudra exécuter **"go run server"** afin que le serveur se lance et attende les connexions des utilisateurs.  

### Lancer le client
Pour lancer le client en ligne de commande, il faut se placer dans le dossier contenant le fichier **client.go**.  
Ensuite, il faudra exécuter **"go run client"** afin que le client se lance et essaye de se connecter au serveur.  
Si aucun serveur n'a été lancé, le client se ferme tout seul.  

Au lancement, il va falloir donner un nom d'utilisateur unique.  
Si le nom que l'on souhaite utiliser est déjà attribué à quelqu'un, le serveur nous le dite et nous demande de sélectionner un autre nom.  

En lançant plusieurs fois la commande **"go run client"**, on peut créer plusieurs clients.  

### Commandes utilisateur
Attention les commandes utilisateurs sont sensible à la casse.  
Il faut aussi faire attention à ne pas mettre un espace au début, sinon la commende ne sera pas reconnue.   
En cas d'erreur dans la saisie d'une commande, la commande **help** est appelée automatiquement pour indiquer les commandes existantes et leur orthographe.  

#### help
Cette commande permet d'afficher un aide-mémoire de toutes les commandes disponibles et de leurs fonctions.  

#### list
Cette commande permet d'afficher la liste de tous les lots actuellement mis aux enchères avec leurs informations de base.  

#### add
Cette commande permet d'ajouter un nouveau lot aux enchères et s'exécute en 3 temps.  
Les 3 actions nécessaires sont :  
- Insérer [name] pour définir le nom du lot  
- Insérer [minPrice] pour définir le prix initial du lot  
- Insérer [duration] au format XhYmZs _(Xh, Ym ou Zs suffisent si l'on souhaite être moins précis)_ pour définir la durée de l'enchère 

#### select [auctionId]
Cette commande permet d'avoir les informations précises d'un lot grâce à son numéro unique [auctionId].  

#### raise [auctionId] [newPrice]
Cette commande permet d'enchérir sur un lot en indiquant le montant que l'on souhaite mettre pour acquérir ce lot.  
Le premier paramètre [auctionId] permet d'identifier le lot sur lequel on souhaite enchérir tandis que le second [newPrice] indique le montant auquel on souhaite enchérir.  

#### addNotifyAllNew
Cette commande permet de recevoir une notification pour chaque nouveau lot créé.  

#### removeNotifyAllNew
Cette commande permet d'arrêter de recevoir une notification pour chaque nouveau lot créé.  

#### addNotify [auctionId]
Cette commande permet de recevoir une notification pour chaque nouvelle enchère sur le lot numéroté [auctionId].  
Par défaut, le propriétaire d'un lot ne reçoit pas de notification à chaque nouvelle enchère contrairement à ce qui est demandé dans la consigne. On s'est rendu compte de ce point un peu trop tardivement.  

#### removeNotify [auctionId]
Cette commande permet d'arrêter de recevoir une notification pour chaque nouvelle enchère sur le lot numéroté [auctionId].  

#### quit
Taper la commande "quit" permet de se déconnecter.  
C'est la seule commande qui fonctionne aussi au lancement, lors de la saisie du nom d'utilisateur.  
On ne peut donc pas se nommer "quit" ;)
 